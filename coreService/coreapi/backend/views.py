from django.shortcuts import render
from .models import CustomUser
from .serializers import CustomUserSerializer
from django.http import HttpRequest
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework import status, viewsets
from .services import model, generate_problem, gen_content, gen_video, chat_generate_script, uploadToS3AndDelete, DummyRequest, SignedUrlStore
import requests
import os
from json import dumps
from uuid import uuid4
from drf_yasg.views import get_schema_view
from drf_yasg import openapi
from django.contrib.auth import authenticate
from rest_framework.authtoken.models import Token
import numpy as np
import chromadb


store = SignedUrlStore()

class CustomUserViewSet(viewsets.ModelViewSet):
    queryset = CustomUser.objects.all()
    serializer_class = CustomUserSerializer

class LoginView(APIView):
    def post(self, request, format=None):
        username = request.data.get("username")
        password = request.data.get("password")
        user = authenticate(username=username, password=password)
        if user is not None:
            token, created = Token.objects.get_or_create(user=user)
            return Response({"token": token.key})
        return Response({"error": "Wrong Credentials"}, status=status.HTTP_400_BAD_REQUEST)


class EmbeddingView(APIView):
    def get(self, request):
        return Response("Make a POST request to this endpoint to get embeddings")

    def post(self, request):
        body = request.data
        embeddings = []
        print(f"docLen: {len(body.get('docs', []))}")
        for i, doc in enumerate(body.get("docs", [])):
            print(i)
            embeddings.append(model.encode(doc, normalize_embeddings=True).tolist())
        print("Done Embedding")
        return Response({
            "embeddings": embeddings
        })

class ReccomendationView(APIView):
    def get(self, request):
        return Response("Make a POST request to this endpoint to get embeddings")

    def post(self, request):
        timeFactor = 0.8
        body = request.data
        docs = body.get('docs', [])
        docs = docs[::-1]
        outVec = model.encode(docs[0], normalize_embeddings=True)
        normalizedLen = 1
        for i, doc in enumerate(docs[1:min(500, len(docs))]):
            outVec += model.encode(doc) * timeFactor
            normalizedLen += timeFactor
            timeFactor *= timeFactor
        outVec = outVec / timeFactor
        res = requests.post(os.getenv("VECTOR_BASE_URL") + "/get_raw", data=dumps({
            "queries": [docs[0]],
            "vectors": [outVec.tolist()],
            "subject": body.get("subject", ""),
            "textbook": body.get("textbook", "")
        }))

        url = res.json()["url"].split("/")
        uuid = url[len(url) - 1]

        return Response({
            "status": res.json()["res"],
            "url": request.build_absolute_uri(f"/query_cache/status/get/{uuid}")
        })

class ProblemGeneratorView(APIView):
    def get(self, request, subject):
        return Response(generate_problem(subject))
    
class VideoGeneratorView(APIView):
    def post(self, request):
        body = request.data
        subject = body.get('subject')
        textbook = body.get('textbook')
        ctx = body.get('textbook_context')
        print(ctx)
        try:
            script = chat_generate_script(subject, ctx)
            uuid = str(uuid4())
            video = gen_video(gen_content(script, uuid), uuid)
            s3Uri = uploadToS3AndDelete(video)
            JITQueryCreateView().post(DummyRequest({
                "docs": [script],
                "subject": subject,
                "textbook": textbook,
                "problems": [{"question": "", "answer": ""}],
                "s3VideoUri": [s3Uri]
            }, request))
            return Response("Video Complete")
        except Exception as e:
            print(f"AWS content blocked {e}")
            return Response("AWS content blocked")
    
class JITQueryStatusCreateView(APIView):
    def get(self, request, uuid):
        return Response(requests.get(f"{os.getenv('VECTOR_BASE_URL')}/status/add/{uuid}").json())

class JITQueryCreateView(APIView):
    def get(self, request):
        return Response("Make a post to this endpoint to cache a query")

    def post(self, request):
        body = request.data
        res = requests.post(os.getenv("VECTOR_BASE_URL") + "/add", data=dumps({
            "subject": body.get("subject", "subject"), 
            "queries": body.get("docs", []),
            "problems": body.get("problems", []),
            "s3VideoUri": body.get("s3VideoUri", []),
            "textbook": body.get("textbook", "")
        }))

        url = res.json()["url"].split("/")
        uuid = url[len(url) - 1]
        print(f"shit {url}")
        return Response({
            "status": res.json()["res"],
            "url": request.build_absolute_uri(f"/query_cache/status/add/{uuid}")
        })

class JITQuerySearchView(APIView):
    def get(self, request):
        return Response("Make a post to compare a query")
    
    def post(self, request):
        body = request.data
        res = requests.post(os.getenv("VECTOR_BASE_URL") + "/get", data=dumps({
            "queries": body["docs"],
            "subject": body.get("subject", ""),
            "textbook": body.get("textbook", "")
        }))

        url = res.json()["url"].split("/")
        uuid = url[len(url) - 1]

        return Response({
            "status": res.json()["res"],
            "url": request.build_absolute_uri(f"/query_cache/status/get/{uuid}")
        })

class JITQueryStatusSearchView(APIView):
    def get(self, request, uuid):
        res = requests.get(f"{os.getenv('VECTOR_BASE_URL')}/status/get/{uuid}").json()
        print(res)
        if (res["status"] == "Finished"):
            for k, v in res["body"].items():
                tmp = []
                for item in v:
                    item["metadata"]["s3VideoUri"] = store.get_signed_url(item["metadata"]["s3VideoUri"])
                    tmp.append(item)
                res["body"][k] = tmp
            
        return Response(res)

class ListSubjectsView(APIView):
    def get(self, request):
        client = chromadb.HttpClient(host='edutok-chroma-1', port = 8000)
        out = dict()
        collections = [i.name for i in client.list_collections()]
        out["collections"] = [collection for collection in collections if "jit" not in collection]
        return Response(out)