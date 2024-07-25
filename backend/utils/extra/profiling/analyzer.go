package profiling

//import (
//	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
//	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
//	"github.com/pyroscope-io/client/pyroscope"
//	log "github.com/sirupsen/logrus"
//	"gorm.io/plugin/opentelemetry/logging/logrus"
//	"os"
//	"runtime"
//)
//
//func InitPyroscope(appName string) {
//	if config.Conf.PyroScope.State != "enable" {
//		logging.Logger.WithFields(log.Fields{
//			"appName": appName,
//		}).Infof("User close Pyroscope, the service would not run.")
//		return
//	}
//
//	runtime.SetMutexProfileFraction(5)
//	runtime.SetBlockProfileRate(5)
//
//	_, err := pyroscope.Start(pyroscope.Config{
//		ApplicationName: appName,
//		ServerAddress:   config.Conf.PyroScope.Addr,
//		Logger:          logrus.NewWriter(),
//		Tags:            map[string]string{"hostname": os.Getenv("HOSTNAME")},
//		ProfileTypes: []pyroscope.ProfileType{
//			pyroscope.ProfileCPU,
//			pyroscope.ProfileAllocObjects,
//			pyroscope.ProfileAllocSpace,
//			pyroscope.ProfileInuseObjects,
//			pyroscope.ProfileInuseSpace,
//			pyroscope.ProfileGoroutines,
//			pyroscope.ProfileMutexCount,
//			pyroscope.ProfileMutexDuration,
//			pyroscope.ProfileBlockCount,
//			pyroscope.ProfileBlockDuration,
//		},
//	})
//
//	if err != nil {
//		logging.Logger.WithFields(log.Fields{
//			"appName": appName,
//			"err":     err,
//		}).Warnf("Pyroscope failed to run.")
//		return
//	}
//}
