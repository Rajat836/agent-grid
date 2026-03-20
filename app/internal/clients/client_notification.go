package clients

import (
	"fmt"

	"bitbucket.org/fyscal/be-commons/pkg/telemetry"
	notificationService "bitbucket.org/fyscal/be-proto/go-proto/notification/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientNotification struct {
	access *clientAccess
	conn   *grpc.ClientConn
	client notificationService.NotificationServiceClient
}

type ClientNotificationMethods interface {
	GetClient() notificationService.NotificationServiceClient
	Close() error
}

func NewClientNotification(access *clientAccess) (*ClientNotification, error) {
	// Get notification service address from config
	notificationAddr := access.cfg.Notification.GrpcHost
	if notificationAddr == "" {
		return nil, fmt.Errorf("notification service address not configured")
	}

	// Create gRPC connection
	// Used telemetry.NewGrpcClient wrapper to get the gRPC client with telemetry enabled
	conn, err := telemetry.NewGrpcClient(
		notificationAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		access.logger.Errorf("Failed to connect to notification service: %v", err)
		return nil, fmt.Errorf("failed to connect to notification service: %w", err)
	}

	client := notificationService.NewNotificationServiceClient(conn)

	access.logger.Infof("Connected to notification service at %s", notificationAddr)

	return &ClientNotification{
		access: access,
		conn:   conn,
		client: client,
	}, nil
}

// GetClient returns the underlying gRPC client
func (c *ClientNotification) GetClient() notificationService.NotificationServiceClient {
	return c.client
}

// Close closes the gRPC connection
func (c *ClientNotification) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
