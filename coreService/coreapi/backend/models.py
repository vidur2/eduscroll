from django.db import models
from django.contrib.auth.models import AbstractUser

class Subject(models.Model):
    name = models.CharField(max_length=100)
    topics = models.JSONField(blank=True, default=list)

    def __str__(self):
        return self.name

class CustomUser(AbstractUser):
    first_name = models.TextField()
    second_name = models.TextField()
    subjects = models.ManyToManyField(Subject, blank=True)
