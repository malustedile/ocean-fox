#!/bin/bash

URL="http://localhost:8080/destinos"

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Expedição Amazônica",
  "categoria": "Brasil",
  "cabines": 7,
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
  "cabines": 8,
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
  "cabines": 6,
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
  "cabines": 5,
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
  "cabines": 4,
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
  "cabines": 9,
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
  "cabines": 7,
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
  "cabines": 6,
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
  "cabines": 8,
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
  "cabines": 12,
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
  "cabines": 5,
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
  "cabines": 6,
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
  "cabines": 7,
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
  "cabines": 5,
  "descricao": {
    "datasDisponiveis": ["2025-03-15", "2025-04-15"],
    "navio": "Dragon Pearl",
    "embarque": "Xangai",
    "desembarque": "Xangai",
    "lugaresVisitados": ["Hong Kong", "Macau", "Taipei"],
    "noites": 7,
    "valorPorPessoa": 5200
  }
}'

curl -X POST $URL -H "Content-Type: application/json" -d '{
  "nome": "Mar do Norte",
  "categoria": "Europa",
  "cabines": 6,
  "descricao": {
    "datasDisponiveis": ["2025-05-20", "2025-06-20"],
    "navio": "Nordic Star",
    "embarque": "Copenhague",
    "desembarque": "Copenhague",
    "lugaresVisitados": ["Estocolmo", "Helsinque", "Tallinn"],
    "noites": 8,
    "valorPorPessoa": 5500
  }
}'
