.PHONY: init run test

init:
	PYTHONPATH=backend python2 -c "from app import init_db; init_db()"

run:
	python2 ./backend/app.py

test: 
	PYTHONPATH=backend python2 ./backend/app_tests.py

