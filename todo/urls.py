from django.urls import path
from . import views

urlpatterns = [
    path('', views.home, name='home'),
    path('todo/<int:todo_id>/edit/', views.edit_todo, name='edit_todo'),
    path('todo/<int:todo_id>/delete/', views.delete_todo, name='delete_todo'),
    path('todo/<int:todo_id>/toggle/', views.toggle_complete, name='toggle_complete'),
]


