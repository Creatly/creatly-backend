FROM mongo
COPY schools.json /schools.json
CMD mongoimport --host mongodb --db coursesTesting --collection schools --type json --file /schools.json --jsonArray