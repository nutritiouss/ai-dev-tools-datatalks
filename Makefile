.PHONY: help install migrate run test clean superuser

help:
	@echo "Available commands:"
	@echo "  make install     - Install dependencies"
	@echo "  make migrate     - Run database migrations"
	@echo "  make run         - Start the development server"
	@echo "  make test        - Run tests"
	@echo "  make superuser   - Create a Django superuser"
	@echo "  make clean       - Remove Python cache files and database"

install:
	pip install -r requirements.txt

migrate:
	python manage.py makemigrations
	python manage.py migrate

run:
	python manage.py runserver

test:
	python manage.py test

superuser:
	python manage.py createsuperuser

clean:
	find . -type d -name "__pycache__" -exec rm -r {} + 2>/dev/null || true
	find . -type f -name "*.pyc" -delete
	find . -type f -name "*.pyo" -delete
	rm -f db.sqlite3


