.PHONY: init run test

init:
	PYTHONPATH=backend python3 -c "from app import init_db; init_db()"

run:
	python3 ./backend/app.py

test: 
	PYTHONPATH=backend python3 ./backend/app_tests.py

