# Django TODO Application

A simple TODO application built with Django as part of the AI Dev Tools Zoomcamp homework.

## Features

- ✅ Create, edit, and delete TODOs
- 📅 Assign due dates to TODOs
- ✓ Mark TODOs as completed/incomplete
- 🎨 Modern, responsive UI
- 🧪 Comprehensive test coverage

## Requirements

- Python 3.8+
- Django 5.0+

## Installation

1. **Install dependencies:**
   ```bash
   pip install -r requirements.txt
   ```
   
   Or using `uv` (recommended):
   ```bash
   uv pip install -r requirements.txt
   ```

2. **Run migrations:**
   ```bash
   python manage.py makemigrations
   python manage.py migrate
   ```

3. **Create a superuser (optional, for admin access):**
   ```bash
   python manage.py createsuperuser
   ```

## Running the Application

Start the development server:
```bash
python manage.py runserver
```

Then open your browser and navigate to: `http://127.0.0.1:8000/`

## Using Makefile

The project includes a Makefile with common commands:

```bash
make install      # Install dependencies
make migrate      # Run migrations
make run          # Start the development server
make test         # Run tests
make superuser    # Create a superuser
make clean        # Clean cache files and database
```

## Running Tests

Run all tests:
```bash
python manage.py test
```

Or using the Makefile:
```bash
make test
```

## Project Structure

```
01-todo/
├── manage.py
├── requirements.txt
├── Makefile
├── README.md
├── todo_project/          # Django project
│   ├── __init__.py
│   ├── settings.py
│   ├── urls.py
│   ├── wsgi.py
│   └── asgi.py
├── todo/                  # Django app
│   ├── migrations/
│   ├── __init__.py
│   ├── admin.py
│   ├── apps.py
│   ├── models.py
│   ├── tests.py
│   ├── views.py
│   └── urls.py
└── templates/             # HTML templates
    ├── base.html
    └── todo/
        ├── home.html
        ├── edit.html
        └── delete.html
```

## Homework Answers

### Question 1: Install Django
**Answer:** `pip install django`

### Question 2: Project and App
**Answer:** `settings.py` - This file contains `INSTALLED_APPS` where you add the app name.

### Question 3: Django Models
**Answer:** Run migrations - After creating models, you need to run `python manage.py makemigrations` and `python manage.py migrate`.

### Question 4: TODO Logic
**Answer:** `views.py` - This is where you implement the view functions that handle the application logic.

### Question 5: Templates
**Answer:** `TEMPLATES['DIRS']` in project's `settings.py` - This is where you configure the directory path for templates.

### Question 6: Tests
**Answer:** `python manage.py test` - This is the command to run Django tests.

## Admin Panel

Access the Django admin panel at: `http://127.0.0.1:8000/admin/`

You'll need to create a superuser first (see Installation step 3).

## Development

This project was created as part of the [AI Dev Tools Zoomcamp](https://github.com/DataTalksClub/ai-dev-tools-zoomcamp) by DataTalksClub.

## License

This project is for educational purposes.


