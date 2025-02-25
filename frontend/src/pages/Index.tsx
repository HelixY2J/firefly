import { motion } from 'framer-motion';
import { ArrowRight } from 'lucide-react';
import { Link } from 'react-router-dom';
import { Navbar } from '../components/Navbar';
import { useEffect, useState } from 'react';

// Generate random positions for fireflies
const generateFireflies = (count: number) => {
  return Array.from({ length: count }, (_, i) => ({
    id: i,
    x: Math.random() * 100,
    y: Math.random() * 100,
    scale: 0.5 + Math.random() * 0.5,
  }));
};

const Index = () => {
  const [fireflies] = useState(() => generateFireflies(15));

  return (
    <div className="min-h-screen relative overflow-hidden">
      <Navbar />
      
      {/* Fireflies */}
      {fireflies.map((firefly) => (
        <motion.div
          key={firefly.id}
          className="absolute w-2 h-2 rounded-full"
          style={{
            left: `${firefly.x}%`,
            top: `${firefly.y}%`,
          }}
          initial={{
            scale: firefly.scale,
            opacity: 0,
          }}
          animate={{
            x: [
              Math.random() * 100 - 50,
              Math.random() * 100 - 50,
              Math.random() * 100 - 50,
            ],
            y: [
              Math.random() * 100 - 50,
              Math.random() * 100 - 50,
              Math.random() * 100 - 50,
            ],
            opacity: [0.2, 0.8, 0.2],
            scale: [firefly.scale, firefly.scale * 1.2, firefly.scale],
          }}
          transition={{
            duration: 8 + Math.random() * 10,
            repeat: Infinity,
            ease: "linear",
          }}
        >
          <motion.div
            className="w-full h-full bg-glow rounded-full"
            animate={{
              boxShadow: [
                "0 0 10px 2px rgba(80, 227, 194, 0.3)",
                "0 0 20px 4px rgba(80, 227, 194, 0.6)",
                "0 0 10px 2px rgba(80, 227, 194, 0.3)",
              ],
            }}
            transition={{
              duration: 2 + Math.random() * 2,
              repeat: Infinity,
              ease: "easeInOut",
            }}
          />
        </motion.div>
      ))}
      
      <main className="container mx-auto px-4 pt-32">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
          className="text-center relative z-10"
        >
          <motion.img
            src="./sparkle.svg"
            alt="Firefly Logo"
            className="w-80 h-80 mx-auto mb-8"
            initial={{ scale: 0.8 }}
            animate={{ scale: 1 }}
            transition={{ duration: 0.5, delay: 0.2 }}
          />
          
          <motion.h1
            className="text-5xl font-bold mb-6 bg-gradient-to-r from-white to-white/60 bg-clip-text text-transparent"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.5, delay: 0.4 }}
          >
            Synchronize Your Music
          </motion.h1>
          
          <motion.p
            className="text-xl text-white/80 mb-12 max-w-2xl mx-auto"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.5, delay: 0.6 }}
          >
            Experience seamless, low-latency music synchronization across your space,
            without the need for wired connections.
          </motion.p>
          
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.5, delay: 0.8 }}
            whileHover={{ scale: 1.05 }}
          >
            <Link
              to="/music"
              className="glass-panel inline-flex items-center px-8 py-4 space-x-2 hover:bg-white/10 transition-colors"
            >
              <span>Get Started</span>
              <ArrowRight size={18} />
            </Link>
          </motion.div>
        </motion.div>
      </main>
    </div>
  );
};

export default Index;