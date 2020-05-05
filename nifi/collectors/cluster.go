package collectors

import (
	"strconv"
	"strings"
	"time"

	"github.com/diarworld/nifi_exporter/nifi/client"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type clusterMetrics struct {
	totalNodes         *prometheus.Desc
	connectedNodes     *prometheus.Desc
	connectingNodes    *prometheus.Desc
	disconnectedNodes  *prometheus.Desc
	disconnectingNodes *prometheus.Desc
	activeThreadCount  *prometheus.Desc
	queuedBytes        *prometheus.Desc
	queuedFlowfiles    *prometheus.Desc
	nodeStartTime      *prometheus.Desc
}

type ClusterCollector struct {
	api *client.Client

	clusterMetrics
}

func NewClusterCollector(api *client.Client, labels map[string]string) *ClusterCollector {
	basicLabels := []string{"node_id"}
	return &ClusterCollector{
		api: api,

		clusterMetrics: clusterMetrics{
			totalNodes: prometheus.NewDesc(
				MetricNamePrefix+"cluster_total_nodes",
				"NiFi cluster total nodes.",
				[]string{},
				labels,
			),
			connectedNodes: prometheus.NewDesc(
				MetricNamePrefix+"cluster_status_connected",
				"NiFi cluster flow controller that is connected to the cluster. A connecting node transitions to connected after the cluster receives the flow controller's first heartbeat. A connected node can transition to disconnecting.",
				[]string{},
				labels,
			),
			connectingNodes: prometheus.NewDesc(
				MetricNamePrefix+"cluster_status_connecting",
				"NiFi cluster flow controller has issued a connection request to the cluster, but has not yet sent a heartbeat. A connecting node can transition to disconnecting or connected. The cluster will not accept any external requests to change the flow while any node is connecting.",
				[]string{},
				labels,
			),
			disconnectedNodes: prometheus.NewDesc(
				MetricNamePrefix+"cluster_status_disconnected",
				"NiFi cluster flow controller that is not connected to the cluster. A disconnected node can transition to connecting.",
				[]string{},
				labels,
			),
			disconnectingNodes: prometheus.NewDesc(
				MetricNamePrefix+"cluster_status_disconnecting",
				"NiFi cluster flow controller that is in the process of disconnecting from the cluster. A disconnecting node will always transition to disconnected.",
				[]string{},
				labels,
			),
			activeThreadCount: prometheus.NewDesc(
				MetricNamePrefix+"cluster_status_active_threads",
				"NiFi cluster node active thread count.",
				append(basicLabels, "node_address"),
				labels,
			),
			queuedBytes: prometheus.NewDesc(
				MetricNamePrefix+"cluster_queued_bytes",
				"NiFi cluster bytes in queue.",
				append(basicLabels, "node_address"),
				labels,
			),
			queuedFlowfiles: prometheus.NewDesc(
				MetricNamePrefix+"cluster_queued_flowfiles",
				"NiFi cluster queued flow files.",
				append(basicLabels, "node_address"),
				labels,
			),
			nodeStartTime: prometheus.NewDesc(
				MetricNamePrefix+"cluster_node_start_time",
				"NiFi cluster node node start time.",
				append(basicLabels, "node_address"),
				labels,
			),
		},
	}
}

func (c *ClusterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.clusterMetrics.totalNodes
	ch <- c.clusterMetrics.connectedNodes
	ch <- c.clusterMetrics.connectingNodes
	ch <- c.clusterMetrics.disconnectedNodes
	ch <- c.clusterMetrics.disconnectingNodes
	ch <- c.clusterMetrics.activeThreadCount
	ch <- c.clusterMetrics.queuedBytes
	ch <- c.clusterMetrics.queuedFlowfiles
	ch <- c.clusterMetrics.nodeStartTime

}

func (c *ClusterCollector) Collect(ch chan<- prometheus.Metric) {
	cluster, err := c.api.GetCluster(true, "")
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.clusterMetrics.totalNodes, err)
		ch <- prometheus.NewInvalidMetric(c.clusterMetrics.connectedNodes, err)
		ch <- prometheus.NewInvalidMetric(c.clusterMetrics.connectingNodes, err)
		ch <- prometheus.NewInvalidMetric(c.clusterMetrics.disconnectedNodes, err)
		ch <- prometheus.NewInvalidMetric(c.clusterMetrics.disconnectingNodes, err)
		ch <- prometheus.NewInvalidMetric(c.clusterMetrics.activeThreadCount, err)
		ch <- prometheus.NewInvalidMetric(c.clusterMetrics.queuedBytes, err)
		ch <- prometheus.NewInvalidMetric(c.clusterMetrics.queuedFlowfiles, err)
		ch <- prometheus.NewInvalidMetric(c.clusterMetrics.nodeStartTime, err)
		return
	}

	connectedNodes := 0
	connectingNodes := 0
	disconnectingNodes := 0
	disconnectedNodes := 0
	if len(cluster.NodeCluster) > 0 {
		for i := range cluster.NodeCluster {
			node := &cluster.NodeCluster[i]

			ch <- prometheus.MustNewConstMetric(
				c.clusterMetrics.activeThreadCount,
				prometheus.GaugeValue,
				float64(node.ActiveThreadCount),
				node.NodeID,
				node.Address,
			)

			layout := "01/02/2006 15:04:05 MSK"
			nodeStartTime, parseError := time.Parse(layout, node.StartTime)
			if parseError == nil {
				ch <- prometheus.MustNewConstMetric(
					c.clusterMetrics.nodeStartTime,
					prometheus.GaugeValue,
					float64(nodeStartTime.Unix()),
					node.NodeID,
					node.Address,
				)
			}

			nodeQueues := strings.Split(node.Queued, " / ")
			nodeQueueFlow, parseError := strconv.ParseUint(nodeQueues[0], 0, 64)
			if parseError == nil {
				ch <- prometheus.MustNewConstMetric(
					c.clusterMetrics.queuedFlowfiles,
					prometheus.GaugeValue,
					float64(nodeQueueFlow),
					node.NodeID,
					node.Address,
				)
			}
			nodeQueuesBytes, parseError := strconv.ParseUint(strings.Split(nodeQueues[1], " ")[0], 0, 64)
			if parseError == nil {
				ch <- prometheus.MustNewConstMetric(
					c.clusterMetrics.queuedBytes,
					prometheus.GaugeValue,
					float64(nodeQueuesBytes),
					node.NodeID,
					node.Address,
				)
			}

			switch {
			case node.Status == "CONNECTED":
				connectedNodes++
			case node.Status == "CONNECTING":
				connectingNodes++
			case node.Status == "CONNECTING":
				disconnectingNodes++
			case node.Status == "DISCONNECTED":
				disconnectedNodes++
			default:
				disconnectedNodes++
			}

			log.WithFields(log.Fields{
				"Address":           node.Address,
				"activeThreadCount": node.ActiveThreadCount,
				"queued":            nodeQueueFlow,
				"queuedBytes":       nodeQueuesBytes,
				"nodeStartTime":     nodeStartTime.Unix(),
				"nodeStatus":        node.Status,
			}).Info("Intialize collect cluster metrics")
		}
		ch <- prometheus.MustNewConstMetric(
			c.clusterMetrics.totalNodes,
			prometheus.GaugeValue,
			float64(len(cluster.NodeCluster)),
		)
		ch <- prometheus.MustNewConstMetric(
			c.clusterMetrics.connectedNodes,
			prometheus.GaugeValue,
			float64(connectedNodes),
		)
		ch <- prometheus.MustNewConstMetric(
			c.clusterMetrics.connectingNodes,
			prometheus.GaugeValue,
			float64(connectingNodes),
		)
		ch <- prometheus.MustNewConstMetric(
			c.clusterMetrics.disconnectedNodes,
			prometheus.GaugeValue,
			float64(disconnectedNodes),
		)
	}

}
