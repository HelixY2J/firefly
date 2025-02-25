import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { Network, AlertCircle } from 'lucide-react';
import { Navbar } from '../components/Navbar';

type NodeData = {
  id: string;
  name: string;
  type: 'server' | 'client';
  x: number;
  y: number;
};

const createNodeObjects = (nodes: string[]): NodeData[] => {
  if (!nodes || nodes.length === 0) return [];

  const serverNode: NodeData = {
    id: 'server',
    name: 'Main Server',
    type: 'server',
    x: 50,
    y: 50
  };

  const clientNodes = nodes.map((hostname, index) => {
    const angle = (2 * Math.PI * index) / nodes.length;
    const radius = 25;
    
    return {
      id: hostname,
      name: hostname,
      type: 'client' as const,
      x: 50 + radius * Math.cos(angle),
      y: 50 + radius * Math.sin(angle)
    };
  });

  return [serverNode, ...clientNodes];
};

const INITIAL_NODES: NodeData[] = [];

const Nodes = () => {
  const [nodes, setNodes] = useState<NodeData[]>(INITIAL_NODES);
  const [hoveredNode, setHoveredNode] = useState<string | null>(null);
  const [wsConnected, setWsConnected] = useState(false);

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8081/ws');

    ws.onopen = () => {
      setWsConnected(true);
      console.log('Connected to nodes WebSocket');
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log('Received nodes data:', data);
        if (data.type === 'nodes_update') {
          setNodes(createNodeObjects(data.nodes || []));
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket Error:', error);
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
            {nodes.length <= 1 ? (
              <div className="flex flex-col items-center justify-center h-full text-center">
                <Network size={48} className="text-glow mb-4" />
                <h3 className="text-xl font-semibold mb-2">No Fireflies Connected</h3>
                <p className="text-gray-400">Run the client application to connect firefly</p>
              </div>
            ) : (
              <>
                {/* Draw lines between nodes */}
                <svg className="absolute inset-0 w-full h-full">
                  {/* Draw lines between all nodes */}
                  {nodes.flatMap((sourceNode, i) => 
                    nodes.slice(i + 1).map((targetNode) => (
                      <motion.line
                        key={`line-${sourceNode.id}-${targetNode.id}`}
                        initial={{ pathLength: 0, opacity: 0 }}
                        animate={{
                          pathLength: 1,
                          opacity: hoveredNode
                            ? (hoveredNode === sourceNode.id || hoveredNode === targetNode.id ? 0.8 : 0.2)
                            : 0.4
                        }}
                        transition={{
                          pathLength: { duration: 0.8, ease: "easeInOut" },
                          opacity: { duration: 0.3 }
                        }}
                        x1={`${sourceNode.x}%`}
                        y1={`${sourceNode.y}%`}
                        x2={`${targetNode.x}%`}
                        y2={`${targetNode.y}%`}
                        stroke="url(#lineGradient)"
                        strokeWidth="2"
                        strokeLinecap="round"
                      />
                    ))
                  )}
                  
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
                      className={`relative ${node.type === 'server' ? 'w-16 h-16' : 'w-12 h-12'} rounded-full flex items-center justify-center ${
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
                      {/* <Network 
                        size={node.type === 'server' ? 24 : 20} 
                        className="text-background z-10" 
                      /> */}
                      <img src='./firefly.svg' alt="Firefly Logo" className="w-10 h-10 z-10" />
                    </motion.div>
                    <div className="absolute top-full mt-2 left-1/2 -translate-x-1/2 whitespace-nowrap">
                      <span className="text-sm font-medium text-white/80">
                        {node.name}
                      </span>
                    </div>
                  </motion.div>
                ))}
              </>
            )}
          </div>
        </div>
      </main>
    </div>
  );
};

export default Nodes;
