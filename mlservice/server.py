
from concurrent import futures
import logging
import threading
import time
from PIL import Image
from transformers import AutoConfig, AutoModel, AutoProcessor
import mlservice_pb2
import mlservice_pb2_grpc
import grpc
import io
import os
import torch
import torch.nn.functional as F

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

MODEL_ID = "google/siglip2-base-patch16-256"
INACTIVITY_TIMEOUT = 300  # 5 minutes in seconds

lock = threading.Lock()

class EncodingServer(mlservice_pb2_grpc.EncodingService):
    def __init__(self) -> None:
        self.device = "cuda" if torch.cuda.is_available() else "cpu"
        self.model = None
        self.processor = AutoProcessor.from_pretrained(MODEL_ID, use_fast=True)
        self.config = AutoConfig.from_pretrained(MODEL_ID)
        self.logit_scale = None
        self._inactivity_timer = None
        self._load_model()

    def _load_model(self):
        """Load the model into memory."""
        if self.model is None:
            self.model = AutoModel.from_pretrained(
                MODEL_ID,
            ).to(self.device)
            self.logit_scale = self.model.logit_scale.item(
            ) if self.model.logit_scale.item() else 4.60517
            logger.info("SigLIP2 model loaded, device: %s", self.device)

    def _unload_model(self):
        """Unload the model from memory after inactivity timeout."""
        with lock:
            if self.model is not None:
                del self.model
                self.model = None
                if self.device == "cuda":
                    torch.cuda.empty_cache()
                logger.info("SigLIP2 model unloaded due to inactivity")

    def _reset_inactivity_timer(self):
        """Reset the inactivity timer. Call this on every request."""
        if self._inactivity_timer is not None:
            self._inactivity_timer.cancel()
        self._inactivity_timer = threading.Timer(INACTIVITY_TIMEOUT, self._unload_model)
        self._inactivity_timer.daemon = True
        self._inactivity_timer.start()

    def _ensure_model_loaded(self):
        """Ensure the model is loaded before processing a request."""
        self._load_model()
        self._reset_inactivity_timer()


    def EncodeText(self, request, context):
        '''
        generate the L2-normalized 768-d embeddings of the texts
        '''
        with lock:
            self._ensure_model_loaded()
            # Enforce max_length padding to match SigLip's training distribution
            inputs = self.processor(
                text=request.input, 
                return_tensors="pt", 
                padding="max_length", 
                max_length=64, # Default for standard SigLip models
                truncation=True
            ).to(self.device)
            
            with torch.no_grad():
                text_embeddings = self.model.get_text_features(**inputs)
                # L2 Normalize along the final dimension (-1)
                normalized_embeddings = F.normalize(text_embeddings, p=2, dim=-1)
                
            output = normalized_embeddings.cpu().flatten().tolist()
            return mlservice_pb2.EncodeResponse(embedding=output)

    def EncodeImage(self, request, context):
        '''
        generate the L2-normalized 768-d embeddings of the images
        '''
        with lock:
            self._ensure_model_loaded()
            image_data = request.image
            image = Image.open(io.BytesIO(image_data))

            inputs = self.processor(
                images=image, return_tensors="pt").to(self.device) # padding=True is unnecessary for images
                
            with torch.no_grad():
                image_embeddings = self.model.get_image_features(**inputs)
                # L2 Normalize along the final dimension (-1)
                normalized_embeddings = F.normalize(image_embeddings, p=2, dim=-1)
                
            output = normalized_embeddings.cpu().flatten().tolist()
            return mlservice_pb2.EncodeResponse(embedding=output)


    def GetImageSpec(self, request, context):
        width = self.config.vision_config.image_size
        height = self.config.vision_config.image_size
        return mlservice_pb2.GetImageSpecResponse(width=width, height=height)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    mlservice_pb2_grpc.add_EncodingServiceServicer_to_server(
        EncodingServer(), server)
    server.add_insecure_port(os.getenv('CLIP_ADDRESS', '0.0.0.0:8173'))
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()
