
import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { Network, AlertCircle } from 'lucide-react';
import { Navbar } from '../components/Navbar';

// Sample hardcoded data
const SAMPLE_NODES = [
  { id: 1, name: 'Main Server', type: 'server', x: 50, y: 50 },
  { id: 2, name: 'Living Room', type: 'client', x: 30, y: 70 },
  { id: 3, name: 'Bedroom', type: 'client', x: 70, y: 70 },
  { id: 4, name: 'Kitchen', type: 'client', x: 40, y: 85 },
  { id: 5, name: 'Office', type: 'client', x: 60, y: 85 },
];

const Nodes = () => {
  const [nodes, setNodes] = useState(SAMPLE_NODES);
  const [hoveredNode, setHoveredNode] = useState<number | null>(null);
  const [wsConnected, setWsConnected] = useState(false);

  useEffect(() => {
    // Create WebSocket connection
    const ws = new WebSocket('ws://your-websocket-url/nodes');

    ws.onopen = () => {
      setWsConnected(true);
      console.log('Connected to nodes WebSocket');
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setNodes(data.nodes);
    };

    ws.onclose = () => {
      setWsConnected(false);
      console.log('Disconnected from nodes WebSocket');
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <div className="min-h-screen">
      <Navbar />
      
      <main className="container mx-auto px-4 pt-32">
        <div className="glass-panel p-8 min-h-[600px] relative">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-semibold flex items-center gap-2">
              <Network className="text-glow" />
              Connected Nodes
            </h2>
            {!wsConnected && (
              <div className="flex items-center gap-2 text-yellow-400">
                <AlertCircle size={16} />
                <span className="text-sm">Demo Mode</span>
              </div>
            )}
          </div>
          
          <div className="relative w-full h-[500px]">
            {/* Draw lines between nodes */}
            <svg className="absolute inset-0 w-full h-full">
              {nodes.map((node) => {
                if (node.type === 'client') {
                  const serverNode = nodes.find(n => n.type === 'server');
                  if (!serverNode) return null;
                  
                  return (
                    <motion.line
                      key={`line-${node.id}`}
                      initial={{ pathLength: 0, opacity: 0 }}
                      animate={{ 
                        pathLength: 1, 
                        opacity: hoveredNode ? (hoveredNode === node.id ? 0.8 : 0.2) : 0.4 
                      }}
                      transition={{ duration: 1.5, repeat: Infinity }}
                      x1={`${serverNode.x}%`}
                      y1={`${serverNode.y}%`}
                      x2={`${node.x}%`}
                      y2={`${node.y}%`}
                      stroke="url(#lineGradient)"
                      strokeWidth="2"
                      strokeLinecap="round"
                    />
                  );
                }
                return null;
              })}
              
              {/* Gradient definition for lines */}
              <defs>
                <linearGradient id="lineGradient" x1="0%" y1="0%" x2="100%" y2="0%">
                  <stop offset="0%" stopColor="#50E3C2" />
                  <stop offset="100%" stopColor="#2D9CDB" />
                </linearGradient>
              </defs>
            </svg>

            {/* Render nodes */}
            {nodes.map((node) => (
              <motion.div
                key={node.id}
                className="absolute"
                style={{
                  left: `${node.x}%`,
                  top: `${node.y}%`,
                  transform: 'translate(-50%, -50%)'
                }}
                initial={{ scale: 0, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                transition={{ type: 'spring', duration: 0.8 }}
                onMouseEnter={() => setHoveredNode(node.id)}
                onMouseLeave={() => setHoveredNode(null)}
              >
                <motion.div
                  className={`relative ${
                    node.type === 'server' ? 'w-16 h-16' : 'w-12 h-12'
                  } rounded-full flex items-center justify-center ${
                    node.type === 'server' ? 'bg-glow' : 'bg-glow/20'
                  }`}
                  animate={{
                    boxShadow: hoveredNode === node.id
                      ? '0 0 20px 4px rgba(80, 227, 194, 0.5)'
                      : '0 0 15px 2px rgba(80, 227, 194, 0.3)'
                  }}
                >
                  <motion.div
                    className="absolute inset-0 rounded-full bg-glow"
                    initial={{ opacity: 0.2 }}
                    animate={{ 
                      opacity: [0.2, 0.4, 0.2],
                      scale: [1, 1.2, 1] 
                    }}
                    transition={{
                      duration: 2,
                      repeat: Infinity,
                      ease: "easeInOut"
                    }}
                    style={{ filter: 'blur(8px)' }}
                  />
                  <Network 
                    size={node.type === 'server' ? 24 : 20} 
                    className="text-background z-10" 
                  />
                </motion.div>
                <div className="absolute top-full mt-2 left-1/2 -translate-x-1/2 whitespace-nowrap">
                  <span className="text-sm font-medium text-white/80">
                    {node.name}
                  </span>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </main>
    </div>
  );
};

export default Nodes;
