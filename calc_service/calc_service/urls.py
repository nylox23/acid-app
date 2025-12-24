from django.urls import path
from calculator.views import calculate_endpoint

urlpatterns = [
    path('calculate', calculate_endpoint),
]