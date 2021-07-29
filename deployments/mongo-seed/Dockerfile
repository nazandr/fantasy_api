FROM mongo

COPY ./players.json /players.json

CMD mongoimport --host mongodb --db fantasy --collection player_cards --file /players.json