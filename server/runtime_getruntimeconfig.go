package server

import (
	"time"

	"github.com/containers/storage/pkg/idtools"
	//"github.com/pkg/errors"
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
	var uidMappings []*pb.LinuxIDMapping
	var gidMappings []*pb.LinuxIDMapping
	if s.defaultIDMappings != nil && !s.defaultIDMappings.Empty() {
		for _, uid := range defaultIDMappings.Uids {
			uidMappings = append(uidMappings, &pb.LinuxIDMapping{
				ContainerId: uint32(uid.ContainerID),
				HostId:      uint32(uid.HostID),
				Size_:       uint32(uid.Size),
			})
		}
		for _, gid := range defaultIDMappings.Gids {
			gidMappings = append(gidMappings, &pb.LinuxIDMapping{
				ContainerId: uint32(gid.ContainerID),
				HostId:      uint32(gid.HostID),
				Size_:       uint32(gid.Size),
			})
		}
	} else {
		uidMappings = append(uidMappings, &pb.LinuxIDMapping{
			ContainerId: uint32(0),
			HostId:      uint32(0),
			Size_:       uint32(4294967295),
		})
		gidMappings = append(gidMappings, &pb.LinuxIDMapping{
			ContainerId: uint32(0),
			HostId:      uint32(0),
			Size_:       uint32(4294967295),
		})
	}
	linuxConfig := &pb.LinuxUserNamespaceConfig{
		UidMappings: uidMappings,
		GidMappings: gidMappings,
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
