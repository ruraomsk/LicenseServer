package custMain

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan CustMess
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan CustMess),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			{
				h.clients[client] = true
			}
		case client := <-h.unregister:
			{
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send)
				}
			}
		case mess := <-h.broadcast:
			{
				for client := range h.clients {
					select {
					case client.send <- mess:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
