curl -i -X POST -d '{
  "color": "white",
  "country": "USA",
  "name": "Number 8",
  "region": "Napa Valley",
  "review": "Great and inexpensive",
  "sweetness": 1,
  "varietal": "Merlot",
  "vineyard": "Asti",
  "vintage": 2012
}' localhost:8080/cellar/accounts/1/bottles
