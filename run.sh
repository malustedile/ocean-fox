sudo fuser -k 3000/tcp
sudo fuser -k 3001/tcp
sudo fuser -k 3002/tcp
sudo fuser -k 3003/tcp
sudo fuser -k 3004/tcp
sudo fuser -k 3005/tcp
sudo fuser -k 5173/tcp

echo "Iniciando Pagamento..."
(cd pagamento-go && go run main.go > ../pagamento.log 2>&1) &
PID_PAGAMENTO=$!
echo "Pagamento iniciado (PID: $PID_PAGAMENTO). Verifique pagamento.log."

echo "Iniciando Marketing..."
(cd marketing-go && go run main.go > ../marketing.log 2>&1) &
PID_MARKETING=$!
echo "Marketing iniciado (PID: $PID_MARKETING). Verifique marketing.log."

echo "Iniciando Bilhete..."
(cd bilhete && go run main.go > ../bilhete.log 2>&1) &
PID_BILHETE=$!
echo "Bilhete iniciado (PID: $PID_BILHETE). Verifique bilhete.log."

echo "Iniciando Reserva..."
(cd reserva-go && go run main.go > ../reserva.log 2>&1) &
PID_RESERVA=$!
echo "Bilhete iniciado (PID: $PID_RESERVA). Verifique reserva.log."

echo "Iniciando Session..."
(cd session-go && go run main.go > ../session.log 2>&1) &
PID_SESSION=$!
echo "Bilhete iniciado (PID: $PID_SESSION). Verifique session.log."

cd frontend && bun run dev &