from django.test import TestCase, Client
from django.urls import reverse
from django.utils import timezone
from datetime import date, timedelta
from .models import Todo


class TodoModelTest(TestCase):
    """Test cases for the Todo model."""
    
    def setUp(self):
        """Set up test data."""
        self.todo = Todo.objects.create(
            title='Test TODO',
            description='This is a test TODO',
            due_date=date.today() + timedelta(days=7)
        )
    
    def test_todo_creation(self):
        """Test that a TODO can be created."""
        self.assertEqual(self.todo.title, 'Test TODO')
        self.assertEqual(self.todo.description, 'This is a test TODO')
        self.assertFalse(self.todo.is_completed)
        self.assertIsNotNone(self.todo.created_at)
    
    def test_todo_str_representation(self):
        """Test the string representation of a TODO."""
        self.assertEqual(str(self.todo), 'Test TODO')
    
    def test_todo_default_completion_status(self):
        """Test that new TODOs are not completed by default."""
        new_todo = Todo.objects.create(title='New TODO')
        self.assertFalse(new_todo.is_completed)
    
    def test_todo_ordering(self):
        """Test that TODOs are ordered by creation date (newest first)."""
        todo1 = Todo.objects.create(title='First TODO')
        todo2 = Todo.objects.create(title='Second TODO')
        todos = list(Todo.objects.all())
        # Newest should be first
        self.assertEqual(todos[0].title, 'Second TODO')
        self.assertEqual(todos[1].title, 'First TODO')


class TodoViewTest(TestCase):
    """Test cases for TODO views."""
    
    def setUp(self):
        """Set up test client and data."""
        self.client = Client()
        self.todo = Todo.objects.create(
            title='Test TODO',
            description='Test description',
            due_date=date.today()
        )
    
    def test_home_view_get(self):
        """Test that home view displays all TODOs."""
        response = self.client.get(reverse('home'))
        self.assertEqual(response.status_code, 200)
        self.assertContains(response, 'Test TODO')
        self.assertTemplateUsed(response, 'todo/home.html')
    
    def test_create_todo(self):
        """Test creating a new TODO."""
        response = self.client.post(reverse('home'), {
            'title': 'New TODO',
            'description': 'New description',
            'due_date': '2025-12-31'
        })
        self.assertEqual(response.status_code, 302)  # Redirect after creation
        self.assertTrue(Todo.objects.filter(title='New TODO').exists())
    
    def test_create_todo_without_title(self):
        """Test that TODO creation fails without a title."""
        initial_count = Todo.objects.count()
        response = self.client.post(reverse('home'), {
            'description': 'No title TODO'
        })
        # Should redirect but not create TODO
        self.assertEqual(Todo.objects.count(), initial_count)
    
    def test_edit_todo_get(self):
        """Test that edit view displays the TODO form."""
        response = self.client.get(reverse('edit_todo', args=[self.todo.id]))
        self.assertEqual(response.status_code, 200)
        self.assertContains(response, self.todo.title)
        self.assertTemplateUsed(response, 'todo/edit.html')
    
    def test_edit_todo_post(self):
        """Test updating a TODO."""
        response = self.client.post(reverse('edit_todo', args=[self.todo.id]), {
            'title': 'Updated TODO',
            'description': 'Updated description',
            'due_date': '2025-12-31'
        })
        self.assertEqual(response.status_code, 302)  # Redirect after update
        self.todo.refresh_from_db()
        self.assertEqual(self.todo.title, 'Updated TODO')
        self.assertEqual(self.todo.description, 'Updated description')
    
    def test_delete_todo_get(self):
        """Test that delete view displays confirmation."""
        response = self.client.get(reverse('delete_todo', args=[self.todo.id]))
        self.assertEqual(response.status_code, 200)
        self.assertContains(response, self.todo.title)
        self.assertTemplateUsed(response, 'todo/delete.html')
    
    def test_delete_todo_post(self):
        """Test deleting a TODO."""
        todo_id = self.todo.id
        response = self.client.post(reverse('delete_todo', args=[todo_id]))
        self.assertEqual(response.status_code, 302)  # Redirect after deletion
        self.assertFalse(Todo.objects.filter(id=todo_id).exists())
    
    def test_toggle_complete(self):
        """Test toggling TODO completion status."""
        initial_status = self.todo.is_completed
        response = self.client.get(reverse('toggle_complete', args=[self.todo.id]))
        self.assertEqual(response.status_code, 302)  # Redirect after toggle
        self.todo.refresh_from_db()
        self.assertNotEqual(self.todo.is_completed, initial_status)
    
    def test_toggle_complete_twice(self):
        """Test that toggling twice returns to original state."""
        initial_status = self.todo.is_completed
        # Toggle once
        self.client.get(reverse('toggle_complete', args=[self.todo.id]))
        self.todo.refresh_from_db()
        self.assertNotEqual(self.todo.is_completed, initial_status)
        # Toggle again
        self.client.get(reverse('toggle_complete', args=[self.todo.id]))
        self.todo.refresh_from_db()
        self.assertEqual(self.todo.is_completed, initial_status)
    
    def test_edit_nonexistent_todo(self):
        """Test that editing a non-existent TODO returns 404."""
        response = self.client.get(reverse('edit_todo', args=[99999]))
        self.assertEqual(response.status_code, 404)
    
    def test_delete_nonexistent_todo(self):
        """Test that deleting a non-existent TODO returns 404."""
        response = self.client.post(reverse('delete_todo', args=[99999]))
        self.assertEqual(response.status_code, 404)


class TodoAdminTest(TestCase):
    """Test cases for Django admin integration."""
    
    def setUp(self):
        """Set up test data."""
        self.todo = Todo.objects.create(
            title='Admin Test TODO',
            description='Test description',
            due_date=date.today()
        )
    
    def test_todo_in_admin(self):
        """Test that Todo model is registered in admin."""
        from django.contrib import admin
        from .models import Todo
        self.assertTrue(admin.site.is_registered(Todo))


