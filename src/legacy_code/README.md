# ¿Who Knows? - Flask application

A search engine from 2009 built with the latest technology! Python 2.7 and Flask 0.5. 

**Note**: This application is intentionally full of problems and vulnerabilities. Do not run it in a production environment. 

## Installation

The dependencies can be installed like this:

```bash
pip install -r requirements.txt
```

But the dependencies are old versions from 2009. An upgrade is recommended:

```bash
pip install -r requirements.txt --upgrade
```

To initialize a new database:

```bash
$ make init
```

Note: Windows does not natively support Make. 


## Running the application

Start a development server on port `8080`:

```bash
$ make run
```
Or:

```bash
$ python2 app.py
```

## Test the application

To run the tests:

```bash
$ make test
```

Or:

```bash
$ python2 ./app_tests.py
```