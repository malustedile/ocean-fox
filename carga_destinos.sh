#!/bin/bash

URL="http://localhost:3000/destinos"

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Expedição Amazônica",
  "categoria": "Brasil",
  "descricao": {
    "datasDisponiveis": ["2025-09-01", "2025-10-15"],
    "navio": "Amazon Star",
    "embarque": "Manaus",
    "desembarque": "Manaus",
    "lugaresVisitados": ["Parintins", "Alter do Chão", "Santarém"],
    "noites": 5,
    "valorPorPessoa": 2800
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Costa da Argentina",
  "categoria": "América do Sul",
  "descricao": {
    "datasDisponiveis": ["2025-01-05", "2025-01-19"],
    "navio": "Southern Light",
    "embarque": "Buenos Aires",
    "desembarque": "Buenos Aires",
    "lugaresVisitados": ["Montevidéu", "Punta del Este", "Puerto Madryn"],
    "noites": 8,
    "valorPorPessoa": 3900
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Caribe Tropical",
  "categoria": "Caribe",
  "descricao": {
    "datasDisponiveis": ["2025-06-10", "2025-07-20"],
    "navio": "Ocean Dream",
    "embarque": "Miami",
    "desembarque": "Miami",
    "lugaresVisitados": ["Bahamas", "Ilhas Cayman", "Jamaica"],
    "noites": 7,
    "valorPorPessoa": 4200
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Alasca Glacial",
  "categoria": "América do Norte",
  "descricao": {
    "datasDisponiveis": ["2025-06-15", "2025-08-01"],
    "navio": "North Explorer",
    "embarque": "Seattle",
    "desembarque": "Seattle",
    "lugaresVisitados": ["Juneau", "Skagway", "Glacier Bay"],
    "noites": 10,
    "valorPorPessoa": 5300
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Costa Oeste Africana",
  "categoria": "África",
  "descricao": {
    "datasDisponiveis": ["2025-11-03", "2025-12-01"],
    "navio": "African Pearl",
    "embarque": "Dacar",
    "desembarque": "Dacar",
    "lugaresVisitados": ["Abidjan", "Lomé", "Cotonou"],
    "noites": 9,
    "valorPorPessoa": 4700
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Mar Vermelho e Arábia",
  "categoria": "Oriente Médio",
  "descricao": {
    "datasDisponiveis": ["2025-03-10", "2025-04-10"],
    "navio": "Desert Queen",
    "embarque": "Dubai",
    "desembarque": "Dubai",
    "lugaresVisitados": ["Mascate", "Salalah", "Aqaba"],
    "noites": 12,
    "valorPorPessoa": 5900
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Encantos do Japão",
  "categoria": "Ásia",
  "descricao": {
    "datasDisponiveis": ["2025-04-01", "2025-05-01"],
    "navio": "Sakura Cruise",
    "embarque": "Tóquio",
    "desembarque": "Tóquio",
    "lugaresVisitados": ["Osaka", "Hiroshima", "Nagasaki"],
    "noites": 7,
    "valorPorPessoa": 5100
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Rivieras Mediterrâneas",
  "categoria": "Mediterrâneo",
  "descricao": {
    "datasDisponiveis": ["2025-05-10", "2025-06-15"],
    "navio": "Mediterranea",
    "embarque": "Barcelona",
    "desembarque": "Barcelona",
    "lugaresVisitados": ["Marselha", "Florença", "Roma"],
    "noites": 8,
    "valorPorPessoa": 4600
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Fiordes da Noruega",
  "categoria": "Escandinávia",
  "descricao": {
    "datasDisponiveis": ["2025-07-01", "2025-08-01"],
    "navio": "Northern Wind",
    "embarque": "Oslo",
    "desembarque": "Oslo",
    "lugaresVisitados": ["Bergen", "Geiranger", "Tromsø"],
    "noites": 10,
    "valorPorPessoa": 6200
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Austrália e Nova Zelândia",
  "categoria": "Oceania",
  "descricao": {
    "datasDisponiveis": ["2025-12-05", "2026-01-10"],
    "navio": "Pacific Breeze",
    "embarque": "Sydney",
    "desembarque": "Sydney",
    "lugaresVisitados": ["Melbourne", "Wellington", "Auckland"],
    "noites": 14,
    "valorPorPessoa": 7200
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Ilhas Gregas",
  "categoria": "Mediterrâneo",
  "descricao": {
    "datasDisponiveis": ["2025-06-05", "2025-07-10"],
    "navio": "Egeu Azul",
    "embarque": "Atenas",
    "desembarque": "Atenas",
    "lugaresVisitados": ["Santorini", "Mykonos", "Rodes"],
    "noites": 9,
    "valorPorPessoa": 4800
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Cruzeiro pelos Emirados",
  "categoria": "Oriente Médio",
  "descricao": {
    "datasDisponiveis": ["2025-11-15", "2025-12-20"],
    "navio": "Dhow Lux",
    "embarque": "Abu Dhabi",
    "desembarque": "Abu Dhabi",
    "lugaresVisitados": ["Dubai", "Khor Fakkan", "Sir Bani Yas"],
    "noites": 6,
    "valorPorPessoa": 4500
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Paraíso Polinésio",
  "categoria": "Oceania",
  "descricao": {
    "datasDisponiveis": ["2025-10-01", "2025-11-15"],
    "navio": "Tiare Dream",
    "embarque": "Papeete",
    "desembarque": "Papeete",
    "lugaresVisitados": ["Bora Bora", "Moorea", "Raiatea"],
    "noites": 10,
    "valorPorPessoa": 6900
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Cruzeiro da Primavera Chinesa",
  "categoria": "Ásia",
  "descricao": {
    "datasDisponiveis": ["2025-03-20", "2025-04-05"],
    "navio": "Red Dragon",
    "embarque": "Xangai",
    "desembarque": "Xangai",
    "lugaresVisitados": ["Pequim", "Qingdao", "Jeju (Coreia)"],
    "noites": 12,
    "valorPorPessoa": 5400
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Cruzeiro da África do Sul",
  "categoria": "África",
  "descricao": {
    "datasDisponiveis": ["2025-12-10", "2026-01-15"],
    "navio": "Cape Spirit",
    "embarque": "Cidade do Cabo",
    "desembarque": "Cidade do Cabo",
    "lugaresVisitados": ["Durban", "Port Elizabeth", "Ilha de Bazaruto"],
    "noites": 11,
    "valorPorPessoa": 6100
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Cruzeiro Sul-Americano de Verão",
  "categoria": "América do Sul",
  "descricao": {
    "datasDisponiveis": ["2025-01-10", "2025-02-10"],
    "navio": "Andes Marinos",
    "embarque": "Santos",
    "desembarque": "Santos",
    "lugaresVisitados": ["Rio de Janeiro", "Búzios", "Montevidéu"],
    "noites": 8,
    "valorPorPessoa": 3700
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Cruzeiro pelas Grandes Antilhas",
  "categoria": "Caribe",
  "descricao": {
    "datasDisponiveis": ["2025-02-15", "2025-03-20"],
    "navio": "Tropical Spirit",
    "embarque": "San Juan",
    "desembarque": "San Juan",
    "lugaresVisitados": ["Cuba", "República Dominicana", "Porto Rico"],
    "noites": 9,
    "valorPorPessoa": 5100
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Cruzeiro pelas Montanhas Rochosas",
  "categoria": "América do Norte",
  "descricao": {
    "datasDisponiveis": ["2025-07-15", "2025-08-20"],
    "navio": "Rocky Ocean",
    "embarque": "Vancouver",
    "desembarque": "Vancouver",
    "lugaresVisitados": ["Victoria", "Prince Rupert", "Sitka"],
    "noites": 10,
    "valorPorPessoa": 5500
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Cruzeiro Brasil Litoral Sul", "categoria": "Brasil", "descricao": {"datasDisponiveis": ["2025-01-10", "2025-02-15"], "navio": "Costa Diadema", "embarque": "Santos", "desembarque": "Santos", "lugaresVisitados": ["Ilhabela", "Balne\u00e1rio Cambori\u00fa", "Porto Belo"], "noites": 7, "valorPorPessoa": 2350}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Ilhas Gregas Encantadas", "categoria": "Mediterr\u00e2neo", "descricao": {"datasDisponiveis": ["2025-06-01", "2025-07-15"], "navio": "MSC Opera", "embarque": "Atenas", "desembarque": "Atenas", "lugaresVisitados": ["Mykonos", "Santorini", "Rodes"], "noites": 10, "valorPorPessoa": 4500}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Caribe Cl\u00e1ssico", "categoria": "Caribe", "descricao": {"datasDisponiveis": ["2025-03-20", "2025-04-10"], "navio": "Harmony of the Seas", "embarque": "Miami", "desembarque": "Miami", "lugaresVisitados": ["Cozumel", "Ilhas Cayman", "Jamaica"], "noites": 7, "valorPorPessoa": 3900}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Expedi\u00e7\u00e3o Patag\u00f4nia Azul", "categoria": "Am\u00e9rica do Sul", "descricao": {"datasDisponiveis": ["2025-11-05", "2025-12-01"], "navio": "Ventus Australis", "embarque": "Ushuaia", "desembarque": "Punta Arenas", "lugaresVisitados": ["Canal de Beagle", "Glaciar Pia", "Ilhas Magalh\u00e3es"], "noites": 6, "valorPorPessoa": 5200}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Navega\u00e7\u00e3o Mediterr\u00e2nea Premium", "categoria": "Mediterr\u00e2neo", "descricao": {"datasDisponiveis": ["2025-08-10", "2025-09-05"], "navio": "Celebrity Edge", "embarque": "Barcelona", "desembarque": "Roma", "lugaresVisitados": ["Marselha", "Floren\u00e7a", "N\u00e1poles"], "noites": 9, "valorPorPessoa": 4800}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Descobertas Norueguesas", "categoria": "Escandin\u00e1via", "descricao": {"datasDisponiveis": ["2025-06-20", "2025-07-25"], "navio": "Norwegian Jade", "embarque": "Oslo", "desembarque": "Copenhague", "lugaresVisitados": ["Bergen", "Geiranger", "Stavanger"], "noites": 8, "valorPorPessoa": 6200}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Cruzeiro Ilhas do Pac\u00edfico", "categoria": "Oceania", "descricao": {"datasDisponiveis": ["2025-02-28", "2025-03-25"], "navio": "Pacific Explorer", "embarque": "Sydney", "desembarque": "Sydney", "lugaresVisitados": ["Ilhas Fiji", "Nova Caled\u00f4nia", "Vanuatu"], "noites": 12, "valorPorPessoa": 5300}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Belezas do Norte da \u00c1frica", "categoria": "\u00c1frica", "descricao": {"datasDisponiveis": ["2025-05-05", "2025-06-01"], "navio": "MSC Musica", "embarque": "Casablanca", "desembarque": "T\u00fanis", "lugaresVisitados": ["T\u00e2nger", "Argel", "La Goulette"], "noites": 7, "valorPorPessoa": 4100}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Tesouros do Oriente M\u00e9dio", "categoria": "Oriente M\u00e9dio", "descricao": {"datasDisponiveis": ["2025-01-15", "2025-03-10"], "navio": "Mein Schiff", "embarque": "Dubai", "desembarque": "Dubai", "lugaresVisitados": ["Abu Dhabi", "Doha", "Muscat"], "noites": 7, "valorPorPessoa": 3700}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Explora\u00e7\u00e3o Sudeste Asi\u00e1tico", "categoria": "\u00c1sia", "descricao": {"datasDisponiveis": ["2025-10-01", "2025-10-20"], "navio": "Quantum of the Seas", "embarque": "Cingapura", "desembarque": "Bangkok", "lugaresVisitados": ["Phuket", "Langkawi", "Ho Chi Minh"], "noites": 9, "valorPorPessoa": 4900}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Costa Brasileira ao Extremo Norte", "categoria": "Brasil", "descricao": {"datasDisponiveis": ["2025-04-05", "2025-05-10"], "navio": "Costa Favolosa", "embarque": "Recife", "desembarque": "Bel\u00e9m", "lugaresVisitados": ["Macei\u00f3", "Natal", "S\u00e3o Lu\u00eds"], "noites": 8, "valorPorPessoa": 3000}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Cruzeiro dos Andes", "categoria": "Am\u00e9rica do Sul", "descricao": {"datasDisponiveis": ["2025-01-20", "2025-03-18"], "navio": "Silver Explorer", "embarque": "Valpara\u00edso", "desembarque": "Buenos Aires", "lugaresVisitados": ["Puerto Montt", "Ushuaia", "Punta del Este"], "noites": 12, "valorPorPessoa": 6200}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Mar do Caribe Sul", "categoria": "Caribe", "descricao": {"datasDisponiveis": ["2025-05-15", "2025-06-10"], "navio": "Caribbean Princess", "embarque": "Cartagena", "desembarque": "Aruba", "lugaresVisitados": ["Bonaire", "Cura\u00e7ao", "Santa Marta"], "noites": 7, "valorPorPessoa": 4400}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Rota da Costa Leste", "categoria": "Am\u00e9rica do Norte", "descricao": {"datasDisponiveis": ["2025-07-05", "2025-07-22"], "navio": "Norwegian Escape", "embarque": "Nova York", "desembarque": "Boston", "lugaresVisitados": ["Halifax", "Portland", "Bar Harbor"], "noites": 9, "valorPorPessoa": 5100}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "\u00c1frica Oriental Selvagem", "categoria": "\u00c1frica", "descricao": {"datasDisponiveis": ["2025-09-10", "2025-10-05"], "navio": "MSC Orchestra", "embarque": "Durban", "desembarque": "Zanzibar", "lugaresVisitados": ["Maputo", "Nosy Be", "Mahajanga"], "noites": 10, "valorPorPessoa": 5500}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Riquezas \u00c1rabes", "categoria": "Oriente M\u00e9dio", "descricao": {"datasDisponiveis": ["2025-02-01", "2025-03-01"], "navio": "MSC Bellissima", "embarque": "Abu Dhabi", "desembarque": "Doha", "lugaresVisitados": ["Sir Bani Yas", "Dammam", "Dubai"], "noites": 7, "valorPorPessoa": 4300}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "J\u00f3ias Asi\u00e1ticas Tropicais", "categoria": "\u00c1sia", "descricao": {"datasDisponiveis": ["2025-08-10", "2025-09-01"], "navio": "Spectrum of the Seas", "embarque": "Hong Kong", "desembarque": "Cingapura", "lugaresVisitados": ["Da Nang", "Nha Trang", "Kota Kinabalu"], "noites": 11, "valorPorPessoa": 5700}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Mediterr\u00e2neo Gourmet", "categoria": "Mediterr\u00e2neo", "descricao": {"datasDisponiveis": ["2025-05-01", "2025-06-01"], "navio": "Azamara Quest", "embarque": "Nice", "desembarque": "Veneza", "lugaresVisitados": ["Cannes", "Cinque Terre", "Split"], "noites": 10, "valorPorPessoa": 6900}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Fiordes Inesquec\u00edveis", "categoria": "Escandin\u00e1via", "descricao": {"datasDisponiveis": ["2025-06-01", "2025-07-01"], "navio": "Hurtigruten Nordlys", "embarque": "Bergen", "desembarque": "Kirkenes", "lugaresVisitados": ["Alesund", "Troms\u00f8", "Alta"], "noites": 9, "valorPorPessoa": 6000}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Aventura Australiana", "categoria": "Oceania", "descricao": {"datasDisponiveis": ["2025-11-01", "2025-11-20"], "navio": "Coral Princess", "embarque": "Brisbane", "desembarque": "Melbourne", "lugaresVisitados": ["Sydney", "Hobart", "Adelaide"], "noites": 12, "valorPorPessoa": 7200}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Encantos do Sul Brasileiro", "categoria": "Brasil", "descricao": {"datasDisponiveis": ["2025-12-15", "2026-01-10"], "navio": "MSC Preziosa", "embarque": "Itaja\u00ed", "desembarque": "Santos", "lugaresVisitados": ["Florian\u00f3polis", "Paranagu\u00e1", "Ilhabela"], "noites": 7, "valorPorPessoa": 3100}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Explorando o Caribe Central", "categoria": "Caribe", "descricao": {"datasDisponiveis": ["2025-10-15", "2025-11-10"], "navio": "Carnival Glory", "embarque": "San Juan", "desembarque": "Miami", "lugaresVisitados": ["St. Maarten", "Dominica", "Ant\u00edgua"], "noites": 8, "valorPorPessoa": 4600}}'

curl -X POST $URL -H "Content-Type: application/json" -d '{"nome": "Viagem \u00c1rtica", "categoria": "Escandin\u00e1via", "descricao": {"datasDisponiveis": ["2025-07-01", "2025-08-01"], "navio": "Havila Capella", "embarque": "Troms\u00f8", "desembarque": "Longyearbyen", "lugaresVisitados": ["Honningsv\u00e5g", "Ny-\u00c5lesund", "Isfjorden"], "noites": 10, "valorPorPessoa": 8800}}'

