package consumers

import (
	"reserva-go/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
    reservaCriadaQueue        = "reserva-criada"
    ReservaCriadaExchange     = "reserva-criada-exc"
    reservaCanceladaQueue     = "reserva-cancelada"
    ReservaCanceladaExchange  = "reserva-cancelada-exc"
    pagamentoAprovadoExchange = "pagamento-aprovado-exc"
    pagamentoAprovadoRK       = "pagamento-aprovado"
    pagamentoRecusadoQueue    = "pagamento-recusado"
    pagamentoRecusadoExchange = "pagamento-recusado-exc"
    pagamentoRecusadoRK       = "pagamento-recusado"
    bilheteGeradoQueue        = "bilhete-gerado"
)

// Função auxiliar para declarar exchange
func declareExchange(ch *amqp.Channel, name, kind string, durable bool) {
    err := ch.ExchangeDeclare(
        name, kind, durable, false, false, false, nil,
    )
    failOnError(err, "Failed to declare exchange '"+name+"'")
}

// Função auxiliar para declarar fila
func declareQueue(ch *amqp.Channel, name string, durable, exclusive bool) amqp.Queue {
    q, err := ch.QueueDeclare(
        name, durable, false, exclusive, false, nil,
    )
    failOnError(err, "Failed to declare queue '"+name+"'")
    return q
}

// Função auxiliar para bind
func bindQueue(ch *amqp.Channel, queue, key, exchange string) {
    err := ch.QueueBind(queue, key, exchange, false, nil)
    failOnError(err, "Failed to bind queue '"+queue+"' to exchange '"+exchange+"'")
}

func InitRabbitMQ() {
    // Exchanges e filas globais
    declareQueue(services.RabbitMQChannelGlobal, reservaCriadaQueue, true, false)
    declareExchange(services.RabbitMQChannelGlobal, ReservaCriadaExchange, "fanout", false)

    // Reserva Cancelada
    chReservaCancelada, err := services.RabbitMQConnection.Channel()
    failOnError(err, "Failed to open channel for reserva cancelada")
    defer chReservaCancelada.Close()
    declareQueue(chReservaCancelada, reservaCanceladaQueue, true, false)
    declareExchange(services.RabbitMQChannelGlobal, ReservaCanceladaExchange, "fanout", false)
    bindQueue(chReservaCancelada, reservaCanceladaQueue, "", ReservaCanceladaExchange)

    // Pagamento Aprovado
    chPagamentoAprovado, err := services.RabbitMQConnection.Channel()
    failOnError(err, "Failed to open channel for pagamento aprovado")
    declareExchange(chPagamentoAprovado, pagamentoAprovadoExchange, "direct", true)
    qPagamentoAprovado := declareQueue(chPagamentoAprovado, "", true, true)
    bindQueue(chPagamentoAprovado, qPagamentoAprovado.Name, pagamentoAprovadoRK, pagamentoAprovadoExchange)
    go consumePagamentoAprovado(chPagamentoAprovado, qPagamentoAprovado.Name)

    // Pagamento Recusado
    chPagamentoRecusado, err := services.RabbitMQConnection.Channel()
    failOnError(err, "Failed to open channel for pagamento recusado")
    declareExchange(chPagamentoRecusado, pagamentoRecusadoExchange, "direct", true)
    qPagamentoRecusado := declareQueue(chPagamentoRecusado, pagamentoRecusadoQueue, true, false)
    bindQueue(chPagamentoRecusado, qPagamentoRecusado.Name, pagamentoRecusadoRK, pagamentoRecusadoExchange)
    go consumePagamentoRecusado(chPagamentoRecusado, qPagamentoRecusado.Name)

    // Promoções
    chPromocoes, err := services.RabbitMQConnection.Channel()
    failOnError(err, "Failed to open channel for promocoes")
    qPromocoes := declareQueue(chPromocoes, "promocoes", true, false)
    go consumePromocoes(chPromocoes, qPromocoes.Name)

    // Bilhete Gerado
    chBilheteGerado, err := services.RabbitMQConnection.Channel()
    failOnError(err, "Failed to open channel for bilhete gerado")
    defer chBilheteGerado.Close()
    declareQueue(chBilheteGerado, bilheteGeradoQueue, true, false)
}