package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
	"github.com/moby/moby/client"
)

type Client struct {
	cli *client.Client
}

type Container struct {
	ID      string
	Name    string
	Image   string
	Status  string
	Running bool
	Ports   string
}

type ContainerStats struct {
	Name    string
	CPU     float64
	Memory  float64
	MemUsed string
	MemLimit string
}

type ContainerInspect struct {
	Name   string
	Status string
	Image  string
	IP 	   string
	Ports []string
	Volumes []string
	Env []string
}

type ContainerOverview struct {
	ID       string
	Name     string
	Image    string
	Status   string
	Running  bool
	Ports    string
	CPU      float64
	Memory   float64
	MemUsed  string
	MemLimit string
}

type ContainerProcess struct {
	PID  string
	Name string
	CPU  string
	Mem  string
}

type PruneSummary struct {
	ContainersDeleted int
	ImagesDeleted     int
	VolumesDeleted    int
	SpaceReclaimed    uint64
}


func New() (*Client, error) {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar con Docker: %w", err)
	}
	return &Client{cli: cli}, nil
}

func (c *Client) List() ([]Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.cli.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("error listando contenedores: %w", err)
	}

	containers := make([]Container, 0, len(result.Items))
seen := make(map[string]bool)

for _, r := range result.Items {
    if seen[r.ID] {
        continue
    }
    seen[r.ID] = true

    name := r.Names[0]
    if len(name) > 0 && name[0] == '/' {
        name = name[1:]
    }

    seenPorts := make(map[string]bool)
    ports := ""
    for _, p := range r.Ports {
        if p.PublicPort != 0 {
            entry := fmt.Sprintf("%d->%d", p.PublicPort, p.PrivatePort)
            if !seenPorts[entry] {
                seenPorts[entry] = true
                if ports != "" {
                    ports += ", "
                }
                ports += entry
            }
        }
    }

    containers = append(containers, Container{
        ID:      r.ID[:12],
        Name:    name,
        Image:   r.Image,
        Status:  r.Status,
        Running: r.State == "running",
        Ports:   ports,
    	})
	}
	
	return containers, nil
}

func (c *Client) Stop(nameOrID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.cli.ContainerStop(ctx, nameOrID, client.ContainerStopOptions{})
	if err != nil {
		return fmt.Errorf("error parando %s: %w", nameOrID, err)
	}
	return nil
}

func (c *Client) Start(nameOrID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.cli.ContainerStart(ctx, nameOrID, client.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("error arrancando %s: %w", nameOrID, err)
	}
	return nil
}

func (c *Client) Restart(nameOrID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := c.cli.ContainerRestart(ctx, nameOrID, client.ContainerRestartOptions{})
	if err != nil {
		return fmt.Errorf("error reiniciando %s: %w", nameOrID, err)
	}
	return nil
}

func (c *Client) Logs(nameOrID string, tail string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := c.cli.ContainerLogs(ctx, nameOrID, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
	})
	if err != nil {
		return fmt.Errorf("error leyendo logs de %s: %w", nameOrID, err)
	}
	defer result.Close()

	io.Copy(os.Stdout, result)
	return nil
}

func (c *Client) Stats() ([]ContainerStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := c.cli.ContainerList(ctx, client.ContainerListOptions{All: false})
	if err != nil {
		return nil, fmt.Errorf("error listando contenedores: %w", err)
	}

	stats := make([]ContainerStats, 0, len(result.Items))

	for _, r := range result.Items {
		s, err := c.cli.ContainerStats(ctx, r.ID, client.ContainerStatsOptions{Stream: false})
		if err != nil {
			continue
		}
		defer s.Body.Close()

		// Decodificamos el JSON que devuelve Docker a un map genérico
		// porque el struct oficial cambió mucho entre versiones
		var data map[string]interface{}
		if err := json.NewDecoder(s.Body).Decode(&data); err != nil {
			continue
		}

		name := r.Names[0]
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
		}

		cpuPercent := calcCPU(data)
		memUsed, memLimit, memPercent := calcMem(data)

		stats = append(stats, ContainerStats{
			Name:     name,
			CPU:      cpuPercent,
			Memory:   memPercent,
			MemUsed:  formatBytes(memUsed),
			MemLimit: formatBytes(memLimit),
		})
	}
	return stats, nil
}

func calcCPU(data map[string]interface{}) float64 {
	defer func() { recover() }()

	cpuStats := data["cpu_stats"].(map[string]interface{})
	preCPU := data["precpu_stats"].(map[string]interface{})

	cpuUsage := cpuStats["cpu_usage"].(map[string]interface{})
	preCPUUsage := preCPU["cpu_usage"].(map[string]interface{})

	totalUsage := cpuUsage["total_usage"].(float64)
	preTotalUsage := preCPUUsage["total_usage"].(float64)
	systemUsage := cpuStats["system_cpu_usage"].(float64)

	// precpu_stats no tiene system_cpu_usage en snapshot único — usamos 0
	preSystemUsage := 0.0
	if v, ok := preCPU["system_cpu_usage"]; ok {
		preSystemUsage = v.(float64)
	}

	cpuDelta := totalUsage - preTotalUsage
	systemDelta := systemUsage - preSystemUsage
	numCPUs := cpuStats["online_cpus"].(float64)

	if systemDelta > 0 && cpuDelta > 0 {
		return (cpuDelta / systemDelta) * numCPUs * 100.0
	}
	return 0
}

func calcMem(data map[string]interface{}) (uint64, uint64, float64) {
	defer func() { recover() }()

	memStats := data["memory_stats"].(map[string]interface{})
	usage := uint64(memStats["usage"].(float64))
	limit := uint64(memStats["limit"].(float64))

	// en kernels nuevos el cache se llama "file" dentro de stats
	cache := uint64(0)
	if stats, ok := memStats["stats"].(map[string]interface{}); ok {
		if v, ok := stats["file"]; ok {
			cache = uint64(v.(float64))
		}
	}

	realUsage := usage - cache
	percent := float64(realUsage) / float64(limit) * 100.0
	return realUsage, limit, percent
}

func formatBytes(b uint64) string {
	mb := float64(b) / 1024 / 1024
	if mb >= 1024 {
		return fmt.Sprintf("%.2f GB", mb/1024)
	}
	return fmt.Sprintf("%.2f MB", mb)
}

func (c *Client) Inspect(nameOrID string) (*ContainerInspect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := c.cli.ContainerInspect(ctx, nameOrID, client.ContainerInspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("error inspeccionando %s: %w", nameOrID, err)
	}

	ip := ""
	for _, net := range data.Container.NetworkSettings.Networks {
		ip = net.IPAddress.String()
		break
	}

	ports := []string{}
	for port, bindings := range data.Container.NetworkSettings.Ports {
		for _, b := range bindings {
			ports = append(ports, fmt.Sprintf("%s->%s", b.HostPort, port))
		}
	}

	volumes := []string{}
	for _, m := range data.Container.Mounts {
		volumes = append(volumes, fmt.Sprintf("%s -> %s", m.Source, m.Destination))
	}

	env := []string{}
	for _, e := range data.Container.Config.Env {
		env = append(env, e)
	}

	return &ContainerInspect{
		Name:    data.Container.Name[1:],
		Status: string(data.Container.State.Status),
		Image:   data.Container.Config.Image,
		IP:      ip,
		Ports:   ports,
		Volumes: volumes,
		Env:     env,
	}, nil
}

func (c *Client) Overview() ([]ContainerOverview, error) {
	containers, err := c.List()
	if err != nil {
		return nil, err
	}

	// Mapa de stats por nombre para cruzar con la lista
	statsMap := make(map[string]ContainerStats)
	stats, _ := c.Stats() // si falla stats, igual mostramos los contenedores
	for _, s := range stats {
		statsMap[s.Name] = s
	}

	overview := make([]ContainerOverview, 0, len(containers))
	for _, cont := range containers {
		s := statsMap[cont.Name]
		overview = append(overview, ContainerOverview{
			ID:       cont.ID,
			Name:     cont.Name,
			Image:    cont.Image,
			Status:   cont.Status,
			Running:  cont.Running,
			Ports:    cont.Ports,
			CPU:      s.CPU,
			Memory:   s.Memory,
			MemUsed:  s.MemUsed,
			MemLimit: s.MemLimit,
		})
	}
	return overview, nil
}

func (c *Client) Top(nameOrID string) ([]ContainerProcess, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.cli.ContainerTop(ctx, nameOrID, client.ContainerTopOptions{})
	if err != nil {
		return nil, fmt.Errorf("error obteniendo procesos de %s: %w", nameOrID, err)
	}

	// result.Titles tiene los headers, result.Processes tiene las filas
	// buscamos los índices de las columnas que nos interesan
	pidIdx, nameIdx, cpuIdx, memIdx := -1, -1, -1, -1
	for i, t := range result.Titles {
		switch t {
		case "PID":
			pidIdx = i
		case "CMD", "COMMAND":
			nameIdx = i
		case "%CPU":
			cpuIdx = i
		case "%MEM":
			memIdx = i
		}
	}

	processes := make([]ContainerProcess, 0, len(result.Processes))
	for _, row := range result.Processes {
		p := ContainerProcess{}
		if pidIdx >= 0 && pidIdx < len(row) {
			p.PID = row[pidIdx]
		}
		if nameIdx >= 0 && nameIdx < len(row) {
			p.Name = row[nameIdx]
		}
		if cpuIdx >= 0 && cpuIdx < len(row) {
			p.CPU = row[cpuIdx]
		}
		if memIdx >= 0 && memIdx < len(row) {
			p.Mem = row[memIdx]
		}
		processes = append(processes, p)
	}
	return processes, nil
}

func (c *Client) Prune() (*PruneSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	summary := &PruneSummary{}

	cr, err := c.cli.ContainerPrune(ctx, client.ContainerPruneOptions{})
	if err != nil {
		return nil, fmt.Errorf("Error purgando contenedores: %w", err)
	}
	summary.ContainersDeleted = len(cr.Report.ContainersDeleted)
	summary.SpaceReclaimed += cr.Report.SpaceReclaimed

	ir, err := c.cli.ImagePrune(ctx, client.ImagePruneOptions{})
	if err != nil {
		return nil, fmt.Errorf("Error purgando imágenes: %w", err)
	}
	summary.ImagesDeleted = len(ir.Report.ImagesDeleted)
	summary.SpaceReclaimed += ir.Report.SpaceReclaimed

	vr, err := c.cli.VolumePrune(ctx, client.VolumePruneOptions{})
	if err != nil {
		return nil, fmt.Errorf("Error purgando volúmenes: %w", err)
	}
	summary.VolumesDeleted = len(vr.Report.VolumesDeleted)
	summary.SpaceReclaimed += vr.Report.SpaceReclaimed

	return summary, nil
}