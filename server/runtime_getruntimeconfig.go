package server

import (
	"time"

	"github.com/containers/storage/pkg/idtools"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	pb "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
)

// GetRuntimeConfigInfo returns the runtime config.
func (s *Server) GetRuntimeConfigInfo(ctx context.Context, req *pb.GetRuntimeConfigInfoRequest) (resp *pb.GetRuntimeConfigInfoResponse, err error) {
	const operation = "get_runtime_config"
	defer func() {
		recordOperation(operation, time.Now())
		recordError(operation, err)
	}()
	logrus.Debugf("GetRuntimeConfigInfo")
	defaultIDMappings := s.getIDMappingsInfo()
	uidMapping := &pb.LinuxIDMapping{ContainerId: uint32(0)}
	gidMapping := &pb.LinuxIDMapping{ContainerId: uint32(0)}

	if s.defaultIDMappings != nil {
		hostID, size := s.getHostIDAndSizeForContainerID(0, defaultIDMappings.Uids)
		if hostID < 0 || size <= 0 {
			return nil, errors.New("usernamespace mapping is enabled at runtime but could not figure out mapping for container UID 0 ")
		}
		uidMapping.HostId = uint32(hostID)
		uidMapping.Size_ = uint32(size)

		hostID, size = s.getHostIDAndSizeForContainerID(0, defaultIDMappings.Gids)
		if hostID < 0 || size <= 0 {
			return nil, errors.New("usernamespace mapping is enabled at runtime but could not figure out mapping for container GID 0 ")
		}
		gidMapping.HostId = uint32(hostID)
		gidMapping.Size_ = uint32(size)
	}
	logrus.Debugf("GetRuntimeConfigInfo: default hostUID %v, hostGID %v", uidMapping.HostId, gidMapping.HostId)
	linuxConfig := &pb.LinuxUserNamespaceConfig{
		UidMappings: []*pb.LinuxIDMapping{uidMapping},
		GidMappings: []*pb.LinuxIDMapping{gidMapping},
	}
	activeRuntimeConfig := &pb.ActiveRuntimeConfig{UserNamespaceConfig: linuxConfig}
	return &pb.GetRuntimeConfigInfoResponse{RuntimeConfig: activeRuntimeConfig}, nil
}

func (s *Server) getHostIDAndSizeForContainerID(containerID int, IDs []idtools.IDMap) (int, int) {
	for _, IDMap := range IDs {
		if IDMap.ContainerID == containerID {
			return IDMap.HostID, IDMap.Size
		}
	}
	return -1, -1
}
