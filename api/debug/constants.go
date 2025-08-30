package debug

// HTML templates for debug pages

const (
	// PProfIndexTemplate is the HTML template for the pprof index page
	PProfIndexTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>pprof</title>
</head>
<body>
    <h1>pprof Debug Information</h1>
    <p>Available profiles:</p>
    <ul>
        <li><a href="%s/debug/pprof/goroutine?debug=1">goroutine</a> - stack traces of all current goroutines</li>
        <li><a href="%s/debug/pprof/heap?debug=1">heap</a> - a sampling of memory allocations of live objects</li>
        <li><a href="%s/debug/pprof/threadcreate?debug=1">threadcreate</a> - stack traces that led to the creation of new OS threads</li>
        <li><a href="%s/debug/pprof/block?debug=1">block</a> - stack traces that led to blocking on synchronization primitives</li>
        <li><a href="%s/debug/pprof/mutex?debug=1">mutex</a> - stack traces of holders of contended mutexes</li>
        <li><a href="%s/debug/pprof/profile">profile</a> - CPU profile (30 seconds)</li>
        <li><a href="%s/debug/pprof/trace?seconds=5">trace</a> - execution trace (5 seconds)</li>
    </ul>
</body>
</html>`

	// DebugIndexTemplate is the HTML template for the main debug dashboard
	DebugIndexTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>Debug Dashboard - %s</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; 
            margin: 40px; 
            background-color: #f5f5f5; 
        }
        .container { 
            max-width: 1200px; 
            margin: 0 auto; 
            background: white; 
            padding: 30px; 
            border-radius: 10px; 
            box-shadow: 0 2px 10px rgba(0,0,0,0.1); 
        }
        h1 { 
            color: #333; 
            border-bottom: 3px solid #007acc; 
            padding-bottom: 10px; 
        }
        h2 { 
            color: #555; 
            margin-top: 30px; 
            margin-bottom: 15px; 
        }
        h3 { 
            color: #666; 
            margin-top: 20px; 
            margin-bottom: 10px; 
            cursor: pointer;
            border-bottom: 1px solid #eee;
            padding-bottom: 5px;
        }
        h3:hover { 
            color: #007acc; 
        }
        ul { 
            list-style-type: none; 
            padding: 0; 
        }
        li { 
            margin: 10px 0; 
            padding: 10px; 
            background: #f8f9fa; 
            border-left: 4px solid #007acc; 
            border-radius: 4px; 
        }
        .clickable-item { 
            cursor: pointer; 
            transition: background-color 0.2s;
        }
        .clickable-item:hover { 
            background: #e3f2fd; 
        }
        .service-info { 
            background: #e7f3ff; 
            padding: 15px; 
            border-radius: 5px; 
            margin-bottom: 20px; 
            border-left: 4px solid #007acc; 
        }
        .grid { 
            display: grid; 
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); 
            gap: 30px; 
            margin-top: 20px; 
        }
        .card { 
            background: #fff; 
            border: 1px solid #ddd; 
            border-radius: 8px; 
            padding: 20px; 
        }
        .status-badge { 
            display: inline-block; 
            padding: 4px 8px; 
            background: #28a745; 
            color: white; 
            border-radius: 12px; 
            font-size: 12px; 
            font-weight: bold; 
        }
        .data-section {
            margin-top: 15px;
            padding: 15px;
            background: #f8f9fa;
            border-radius: 5px;
            border-left: 4px solid #007acc;
            display: none;
        }
        .data-section.visible {
            display: block;
        }
        .loading {
            color: #666;
            font-style: italic;
        }
        .error {
            color: #dc3545;
            background: #f8d7da;
            border-left-color: #dc3545;
        }
        pre {
            background: #f8f9fa;
            padding: 10px;
            border-radius: 4px;
            overflow-x: auto;
            white-space: pre-wrap;
            font-size: 12px;
        }
        table {
            width: 100%%;
            border-collapse: collapse;
            margin-top: 10px;
        }
        th, td {
            text-align: left;
            padding: 8px;
            border-bottom: 1px solid #ddd;
        }
        th {
            background-color: #f2f2f2;
            font-weight: bold;
        }
        .pprof-links {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
            margin-top: 10px;
        }
        .pprof-link {
            background: #007acc;
            color: white;
            padding: 8px 12px;
            border-radius: 4px;
            text-decoration: none;
            font-size: 12px;
        }
        .pprof-link:hover {
            background: #005999;
            color: white;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Debug Dashboard</h1>
        
        <div class="service-info">
            <strong>Service:</strong> %s &nbsp;
            <strong>Environment:</strong> %s &nbsp;
            <span class="status-badge">RUNNING</span>
        </div>

        <div class="grid">
            <div class="card">
                <h2>Health Checks</h2>
                
                <h3 onclick="fetchData('health', '%s/health')" class="clickable-item">
                    Service Health
                </h3>
                <div id="health-data" class="data-section"></div>
                
                <h3 onclick="fetchData('readiness', '%s/health/ready')" class="clickable-item">
                    Readiness Probe
                </h3>
                <div id="readiness-data" class="data-section"></div>
                
                <h3 onclick="fetchData('liveness', '%s/health/live')" class="clickable-item">
                    Liveness Probe
                </h3>
                <div id="liveness-data" class="data-section"></div>
            </div>

            <div class="card">
                <h2>System Information</h2>
                
                <h3 onclick="fetchData('service-info', '%s/debug/info')" class="clickable-item">
                    Service Info
                </h3>
                <div id="service-info-data" class="data-section"></div>
                
                <h3 onclick="fetchData('runtime', '%s/debug/runtime')" class="clickable-item">
                    Runtime Stats
                </h3>
                <div id="runtime-data" class="data-section"></div>
                
                <h3 onclick="fetchData('memory', '%s/debug/memory')" class="clickable-item">
                    Memory Stats
                </h3>
                <div id="memory-data" class="data-section"></div>
                
                <h3 onclick="fetchData('vars', '%s/debug/vars')" class="clickable-item">
                    Runtime Variables
                </h3>
                <div id="vars-data" class="data-section"></div>
                
                <h3 onclick="fetchData('env', '%s/debug/env')" class="clickable-item">
                    Environment Variables
                </h3>
                <div id="env-data" class="data-section"></div>
            </div>
        </div>

        %s
        
        <div style="margin-top: 40px; padding-top: 20px; border-top: 1px solid #ddd; color: #666; font-size: 14px;">
            <p>Debug endpoints are enabled for development and debugging purposes. 
            <strong>Ensure these are disabled in production environments.</strong></p>
            <p>Generated at: <span id="timestamp"></span></p>
        </div>
    </div>
    
    <script>
        document.getElementById('timestamp').textContent = new Date().toLocaleString();
        
        // Track which sections are currently visible
        const visibleSections = new Set();
        
        async function fetchData(sectionId, url) {
            const dataDiv = document.getElementById(sectionId + '-data');
            
            // Toggle visibility
            if (visibleSections.has(sectionId)) {
                dataDiv.classList.remove('visible');
                visibleSections.delete(sectionId);
                return;
            }
            
            // Show loading state
            dataDiv.innerHTML = '<div class="loading">Loading...</div>';
            dataDiv.classList.add('visible');
            visibleSections.add(sectionId);
            
            try {
                const response = await fetch(url);
                
                if (!response.ok) {
                    throw new Error('HTTP ' + response.status + ': ' + response.statusText);
                }
                
                const contentType = response.headers.get('content-type');
                let data;
                
                if (contentType && contentType.includes('application/json')) {
                    data = await response.json();
                    dataDiv.innerHTML = formatJsonData(data, sectionId);
                } else {
                    data = await response.text();
                    dataDiv.innerHTML = '<pre>' + escapeHtml(data) + '</pre>';
                }
                
            } catch (error) {
                dataDiv.innerHTML = '<div class="error">Error loading data: ' + escapeHtml(error.message) + '</div>';
                dataDiv.classList.add('error');
            }
        }
        
        function formatJsonData(data, sectionId) {
            if (sectionId === 'memory' || sectionId === 'runtime' || sectionId === 'service-info') {
                return formatAsTable(data);
            } else if (sectionId === 'health' || sectionId === 'readiness' || sectionId === 'liveness') {
                return formatHealthData(data);
            } else {
                return '<pre>' + JSON.stringify(data, null, 2) + '</pre>';
            }
        }
        
        function formatAsTable(data) {
            let html = '<table>';
            html += '<thead><tr><th>Property</th><th>Value</th></tr></thead>';
            html += '<tbody>';
            
            for (const [key, value] of Object.entries(data)) {
                let displayValue = value;
                
                // Format specific values
                if (key.includes('alloc') || key.includes('sys')) {
                    displayValue = formatBytes(value);
                } else if (key === 'pause_total_ns') {
                    displayValue = formatNanoseconds(value);
                } else if (typeof value === 'object') {
                    displayValue = JSON.stringify(value);
                }
                
                html += '<tr><td>' + escapeHtml(key.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())) + '</td>';
                html += '<td>' + escapeHtml(displayValue.toString()) + '</td></tr>';
            }
            
            html += '</tbody></table>';
            return html;
        }
        
        function formatHealthData(data) {
            let html = '<table>';
            html += '<thead><tr><th>Property</th><th>Value</th></tr></thead>';
            html += '<tbody>';
            
            for (const [key, value] of Object.entries(data)) {
                let displayValue = value;
                if (key === 'timestamp') {
                    displayValue = new Date(value).toLocaleString();
                }
                
                html += '<tr><td>' + escapeHtml(key.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())) + '</td>';
                html += '<td>';
                
                if (key === 'status') {
                    const statusClass = value === 'ok' || value === 'ready' || value === 'alive' ? 'status-badge' : 'error';
                    html += '<span class="' + statusClass + '">' + escapeHtml(displayValue.toString()) + '</span>';
                } else {
                    html += escapeHtml(displayValue.toString());
                }
                
                html += '</td></tr>';
            }
            
            html += '</tbody></table>';
            return html;
        }
        
        function formatBytes(bytes) {
            if (bytes === 0) return '0 Bytes';
            const k = 1024;
            const sizes = ['Bytes', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }
        
        function formatNanoseconds(ns) {
            if (ns < 1000) return ns + ' ns';
            if (ns < 1000000) return (ns / 1000).toFixed(2) + ' μs';
            if (ns < 1000000000) return (ns / 1000000).toFixed(2) + ' ms';
            return (ns / 1000000000).toFixed(2) + ' s';
        }
        
        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }
    </script>
</body>
</html>`

	// PProfSectionTemplate is the HTML template for the pprof section in the debug dashboard
	PProfSectionTemplate = `
        <div class="card">
            <h2>Profiling & Performance</h2>
            
            <h3 onclick="showPprofLinks()" class="clickable-item">
                Profiling Tools
            </h3>
            <div id="pprof-links" class="data-section">
                <div class="pprof-links">
                    <a href="%s/debug/pprof/" target="_blank" class="pprof-link">pprof Index</a>
                    <a href="%s/debug/pprof/goroutine?debug=1" target="_blank" class="pprof-link">Goroutines</a>
                    <a href="%s/debug/pprof/heap?debug=1" target="_blank" class="pprof-link">Heap</a>
                    <a href="%s/debug/pprof/profile" target="_blank" class="pprof-link">CPU Profile</a>
                    <a href="%s/debug/pprof/trace?seconds=5" target="_blank" class="pprof-link">Execution Trace</a>
                    <a href="%s/debug/pprof/block?debug=1" target="_blank" class="pprof-link">Block Profile</a>
                    <a href="%s/debug/pprof/mutex?debug=1" target="_blank" class="pprof-link">Mutex Profile</a>
                </div>
                <p style="margin-top: 15px; font-size: 12px; color: #666;">
                    <strong>Note:</strong> Profiling links open in new tabs. CPU profiling takes 30 seconds, execution trace takes 5 seconds.
                </p>
            </div>
        </div>
        
        <script>
            function showPprofLinks() {
                const linksDiv = document.getElementById('pprof-links');
                if (linksDiv.classList.contains('visible')) {
                    linksDiv.classList.remove('visible');
                } else {
                    linksDiv.classList.add('visible');
                }
            }
        </script>`
)
