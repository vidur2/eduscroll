"""
URL configuration for coreapi project.

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/4.2/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.contrib import admin
from django.urls import path, re_path
from backend.views import JITQueryCreateView, JITQueryStatusCreateView, EmbeddingView, JITQuerySearchView, JITQueryStatusSearchView, ProblemGeneratorView, VideoGeneratorView, ReccomendationView, ListSubjectsView
from drf_yasg.views import get_schema_view
from drf_yasg import openapi


schema_view = get_schema_view(
    openapi.Info(
        title="EduScroll Core",
        default_version='v1',
        description="API documentation",
    ),
    public=True, # prolly should change thse for prod
)



urlpatterns = [
    re_path(r'^swagger(?P<format>\.json|\.yaml)$', schema_view.without_ui(cache_timeout=0), name='schema-json'),
    path('swagger/', schema_view.with_ui('swagger', cache_timeout=0), name='schema-swagger-ui'),    
    path('admin/', admin.site.urls),
    path('query_cache/add', JITQueryCreateView.as_view()),
    path('query_cache/search', JITQuerySearchView.as_view()),
    path('query_cache/status/add/<str:uuid>', JITQueryStatusCreateView.as_view()),
    path('query_cache/status/get/<str:uuid>', JITQueryStatusSearchView.as_view()),
    path('embed', EmbeddingView.as_view()),
    path('problem/<str:subject>', ProblemGeneratorView.as_view()),
    path('query_cache/get_reccomendation', ReccomendationView.as_view()),
    path('video', VideoGeneratorView.as_view()),
    path('subjects', ListSubjectsView.as_view()),
]
