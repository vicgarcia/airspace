#!/bin/bash

# make directory for the db setup contents
rm -rf db
mkdir db

# drop into /db
pushd db

# create virtual environmeet for temporary use
echo "setting up environment ..."
# assuming python 3.8 installed as described in the readme
python3.8 -m venv venv
venv/bin/python -m pip install --upgrade pip
cat > requirements.txt << EOF
certifi==2021.5.30
charset-normalizer==2.0.6
faa-aircraft-registry==0.1.0
idna==3.2
requests==2.26.0
urllib3==1.26.7
EOF
venv/bin/pip install -r requirements.txt

# write the python script and execute
cat > setup_faa_sqlite_db.py << EOF
import sqlite3
import json
import requests
import zipfile
import faa_aircraft_registry

# download the FAA aircraft registration database as a zip archive
print('downloading faa aircraft registration data ...')
FAA_DOWNLOAD_URL = 'https://registry.faa.gov/database/ReleasableAircraft.zip'
response = requests.get(FAA_DOWNLOAD_URL)
with open('faa_registry.zip', 'wb') as file:
    file.write(response.content)

# https://www.faa.gov/licenses_certificates/aircraft_certification/aircraft_registry/releasable_aircraft_download/

# unzip the archive and parse
print('parse faa zip archive ...')
data = None
with zipfile.ZipFile('faa_registry.zip') as file:
    data = faa_aircraft_registry.read(file)

# connect to sqlite db and create aircraft table
print('create aircrafts table in sqlite3 db ...')
connection = sqlite3.connect('faa_registry.db')
cursor = connection.cursor()
cursor.execute("create table aircrafts (transponder varchar(6), registration varchar(6), data json)")

# insert the parsed aircrafts into the database
print('insert faa data into aircrafts data ...')
inserted = 0
for row in data['master'].values():
    # print(row['transponder_code_hex'], row['registration_number'])
    cursor.execute(
        'insert into aircrafts (transponder, registration, data) values (?, ?, ?)',
        ( row['transponder_code_hex'], row['registration_number'], json.dumps(row) )
    )
    connection.commit()
    # print output every 20000 rows
    inserted += 1
    if inserted % 20000 == 0:
        print(f'... inserted {inserted} records')

# print total inserted
print(f'insert complete, {inserted} records')

# setup indexes on the aircrafts table
cursor.execute('create index aircrafts_transponder_index on aircrafts(transponder)')
cursor.execute('create index aircrafts_registration_index on aircrafts(registration)')

# close database connection
print('setup complete')
cursor.close()
connection.close()
EOF
venv/bin/python setup_faa_sqlite_db.py

# drop back to the project root
popd

# move the sqlite db into the application path
mv db/faa_registry.db cmd/airspace/faa_registry.db
