sudo fuser -k 3000/tcp
sudo fuser -k 3001/tcp
sudo fuser -k 3002/tcp
sudo fuser -k 3003/tcp
sudo fuser -k 3004/tcp
sudo fuser -k 3005/tcp
sudo fuser -k 5173/tcp
sudo fuser -k 8080/tcp

echo "Iniciando Pagamento..."
(cd backend/pagamento-go && go run main.go > ../../logs/pagamento.log 2>&1) &
PID_PAGAMENTO=$!
echo "Pagamento iniciado (PID: $PID_PAGAMENTO). Verifique pagamento.log."

echo "Iniciando Marketing..."
(cd backend/marketing && go run main.go > ../../logs/marketing.log 2>&1) &
PID_MARKETING=$!
echo "Marketing iniciado (PID: $PID_MARKETING). Verifique marketing.log."

echo "Iniciando Bilhete..."
(cd backend/bilhete && go run main.go > ../../logs/bilhete.log 2>&1) &
PID_BILHETE=$!
echo "Bilhete iniciado (PID: $PID_BILHETE). Verifique bilhete.log."

echo "Iniciando Reserva..."
(cd backend/reserva && go run main.go > ../../logs/reserva.log 2>&1) &
PID_RESERVA=$!
echo "Reserva iniciado (PID: $PID_RESERVA). Verifique reserva.log."

echo "Iniciando Session..."
(cd backend/session && go run main.go > ../../logs/session.log 2>&1) &
PID_SESSION=$!
echo "Session iniciado (PID: $PID_SESSION). Verifique session.log."

echo "Iniciando Itinerarios..."
(cd backend/itinerarios && go run main.go > ../../logs/itinerarios.log 2>&1) &
PID_ITINERARIOS=$!
echo "Itinerarios iniciado (PID: $PID_ITINERARIOS). Verifique itinerarios.log."

cd frontend && bun run dev &