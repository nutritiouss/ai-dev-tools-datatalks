from django.shortcuts import render, redirect, get_object_or_404
from django.contrib import messages
from .models import Todo


def home(request):
    """Display all TODOs and handle creation."""
    todos = Todo.objects.all()
    
    if request.method == 'POST':
        title = request.POST.get('title')
        description = request.POST.get('description', '')
        due_date = request.POST.get('due_date') or None
        
        if title:
            Todo.objects.create(
                title=title,
                description=description,
                due_date=due_date
            )
            messages.success(request, 'TODO created successfully!')
            return redirect('home')
        else:
            messages.error(request, 'Title is required!')
    
    return render(request, 'todo/home.html', {'todos': todos})


def edit_todo(request, todo_id):
    """Edit an existing TODO."""
    todo = get_object_or_404(Todo, id=todo_id)
    
    if request.method == 'POST':
        todo.title = request.POST.get('title', todo.title)
        todo.description = request.POST.get('description', '')
        todo.due_date = request.POST.get('due_date') or None
        
        if todo.title:
            todo.save()
            messages.success(request, 'TODO updated successfully!')
            return redirect('home')
        else:
            messages.error(request, 'Title is required!')
    
    return render(request, 'todo/edit.html', {'todo': todo})


def delete_todo(request, todo_id):
    """Delete a TODO."""
    todo = get_object_or_404(Todo, id=todo_id)
    
    if request.method == 'POST':
        todo.delete()
        messages.success(request, 'TODO deleted successfully!')
        return redirect('home')
    
    return render(request, 'todo/delete.html', {'todo': todo})


def toggle_complete(request, todo_id):
    """Toggle the completion status of a TODO."""
    todo = get_object_or_404(Todo, id=todo_id)
    todo.is_completed = not todo.is_completed
    todo.save()
    messages.success(request, f'TODO marked as {"completed" if todo.is_completed else "incomplete"}!')
    return redirect('home')


