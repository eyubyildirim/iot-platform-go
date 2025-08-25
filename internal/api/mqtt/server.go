package mqtt

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iot-platform/internal/model"
	"iot-platform/internal/service"
	"log"
	"net"
	"strings"
	"time"
)

type Server struct {
	addr              string
	sensorDataService *service.SensorDataService
}

func NewServer(addr string, sensorDataService *service.SensorDataService) *Server {
	return &Server{addr: addr, sensorDataService: sensorDataService}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	log.Printf("MQTT server listening on %s", s.addr)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("accept error: %v", err)
				continue
			}
			go s.handleConn(conn)
		}
	}()
	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	if err := s.handleConnect(r, conn); err != nil {
		log.Printf("connect error: %v", err)
		return
	}
	for {
		if err := s.handlePublish(r); err != nil {
			if err != io.EOF {
				log.Printf("publish error: %v", err)
			}
			return
		}
	}
}

func (s *Server) handleConnect(r *bufio.Reader, conn net.Conn) error {
	header, err := r.ReadByte()
	if err != nil {
		return err
	}
	if header>>4 != 1 {
		return fmt.Errorf("expected CONNECT packet")
	}
	rl, err := r.ReadByte()
	if err != nil {
		return err
	}
	buf := make([]byte, int(rl))
	if _, err := io.ReadFull(r, buf); err != nil {
		return err
	}
	_, err = conn.Write([]byte{0x20, 0x02, 0x00, 0x00})
	return err
}

func (s *Server) handlePublish(r *bufio.Reader) error {
	header, err := r.ReadByte()
	if err != nil {
		return err
	}
	if header>>4 != 3 {
		return fmt.Errorf("expected PUBLISH packet")
	}
	rl, err := r.ReadByte()
	if err != nil {
		return err
	}
	buf := make([]byte, int(rl))
	if _, err := io.ReadFull(r, buf); err != nil {
		return err
	}
	topicLen := int(buf[0])<<8 | int(buf[1])
	topic := string(buf[2 : 2+topicLen])
	payload := buf[2+topicLen:]
	parts := strings.Split(topic, "/")
	if len(parts) != 2 || parts[0] != "sensor-data" {
		return nil
	}
	deviceID := parts[1]
	var msg struct {
		MetricName  string  `json:"metricName"`
		MetricValue float64 `json:"metricValue"`
	}
	if err := json.Unmarshal(payload, &msg); err != nil {
		return err
	}
	data := &model.SensorData{
		DeviceId:    deviceID,
		MetricName:  msg.MetricName,
		MetricValue: msg.MetricValue,
		Timestamp:   time.Now(),
	}
	ctx := context.Background()
	return s.sensorDataService.CreateSensorData(ctx, data)
}
