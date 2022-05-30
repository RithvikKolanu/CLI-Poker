package main

func (s *server) addToPot(c *client, val int) {
	s.pool = s.pool + val
	c.bankroll = c.bankroll - val
}

func (s *server) collect(val int, c *client) {
	c.bankroll = c.bankroll + s.pool
	s.pool = 0
}
