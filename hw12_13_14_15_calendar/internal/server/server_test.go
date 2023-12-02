package server_test

import (
	"context"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	proto "github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/api"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateEvent(t *testing.T) {
	configFile := "../../configs/config.yml"
	cfg, cfgErr := getConfig(t, configFile)
	if cfgErr != nil {
		require.Fail(t, "Test init error", cfgErr)
	}

	ctx := context.Background()

	// Set up a connection to the server.
	conn, connErr := getGRPCConnection(t, ctx, cfg)
	if connErr != nil {
		require.Fail(t, "Test init error", connErr)
	}
	defer conn.Close()

	client := proto.NewEventServiceClient(conn)

	// Successful create
	date := time.Now().AddDate(0, 1, 0).UTC().Round(time.Minute)
	event1Dto := &proto.CreateRequest{EventData: &proto.EventData{
		Title:       "test event 1",
		Description: "test event description",
		DateTime:    timestamppb.New(date),
		RemindTime:  timestamppb.New(date.Add(-time.Hour * 2)),
		Duration:    "01:00:00",
		UserId:      1,
	}}
	resp, err := client.CreateEvent(ctx, event1Dto)
	eventId := resp.EventId.Id

	require.NoError(t, err)
	require.NotEmpty(t, resp.EventId.Id)

	// Date is busy
	event2Dto := &proto.CreateRequest{EventData: &proto.EventData{
		Title:       "test event 2",
		Description: "test event description",
		DateTime:    timestamppb.New(date),
		RemindTime:  timestamppb.New(date.Add(-time.Hour * 1)),
		Duration:    "01:00:00",
		UserId:      1,
	}}
	resp, err = client.CreateEvent(ctx, event2Dto)

	require.NotNil(t, err)
	grpcStatus := status.Convert(err)
	require.NotNil(t, grpcStatus)
	require.Equal(t, codes.Unknown, grpcStatus.Code())
	require.Equal(t, app.ErrDateBusy.Error(), grpcStatus.Message())

	// Remove created
	_, err = client.DeleteEvent(ctx, &proto.DeleteRequest{EventId: &proto.EventId{Id: eventId}})
	require.NoError(t, err)
}

func TestGetEvents(t *testing.T) {
	configFile := "../../configs/config.yml"
	cfg, cfgErr := getConfig(t, configFile)
	if cfgErr != nil {
		require.Fail(t, "Test init error", cfgErr)
	}

	ctx := context.Background()

	// Set up a connection to the server.
	conn, connErr := getGRPCConnection(t, ctx, cfg)
	if connErr != nil {
		require.Fail(t, "Test init error", connErr)
	}
	defer conn.Close()

	client := proto.NewEventServiceClient(conn)

	// Fill db
	dates := []time.Time{
		time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 3, 13, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 10, 13, 0, 0, 0, time.UTC),
	}
	createdEventIds := make([]string, 0, len(dates))

	for i := 0; i < len(dates); i++ {
		dto := &proto.CreateRequest{EventData: &proto.EventData{
			Title:       "test event #" + strconv.Itoa(i),
			Description: "test event description",
			DateTime:    timestamppb.New(dates[i]),
			RemindTime:  timestamppb.New(dates[i].Add(-time.Hour * 2)),
			Duration:    "01:00:00",
			UserId:      1,
		}}
		resp, err := client.CreateEvent(ctx, dto)

		require.NoError(t, err)
		require.NotEmpty(t, resp.EventId.Id)

		createdEventIds = append(createdEventIds, resp.EventId.Id)
	}

	// Test Get*Events routes
	resp, err := client.GetDayEvents(ctx, &proto.StartDate{StartDate: timestamppb.New(dates[0])})
	require.NoError(t, err)
	require.Len(t, resp.GetEvents(), 1)

	resp, err = client.GetWeekEvents(ctx, &proto.StartDate{StartDate: timestamppb.New(dates[0])})
	require.NoError(t, err)
	require.Len(t, resp.GetEvents(), 2)

	resp, err = client.GetMonthEvents(ctx, &proto.StartDate{StartDate: timestamppb.New(dates[0])})
	require.NoError(t, err)
	require.Len(t, resp.GetEvents(), 3)

	// Remove created
	for _, eventId := range createdEventIds {
		_, err := client.DeleteEvent(ctx, &proto.DeleteRequest{EventId: &proto.EventId{Id: eventId}})
		require.NoError(t, err)
	}
}

func getConfig(t *testing.T, configFilePath string) (*config.Config, error) {
	t.Helper()

	file, fileErr := os.Open(configFilePath)
	if fileErr != nil {
		return nil, fileErr
	}

	cfg, configErr := config.New(file)
	if configErr != nil {
		return nil, configErr
	}

	return cfg, nil
}

func getGRPCConnection(t *testing.T, ctx context.Context, config *config.Config) (*grpc.ClientConn, error) {
	t.Helper()

	return grpc.DialContext(
		ctx,
		net.JoinHostPort(config.GRPCServer.Host, config.GRPCServer.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
