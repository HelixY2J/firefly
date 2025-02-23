
import { motion } from 'framer-motion';
import { ArrowRight } from 'lucide-react';
import { Link } from 'react-router-dom';
import { Navbar } from '../components/Navbar';

const Index = () => {
  return (
    <div className="min-h-screen">
      <Navbar />
      
      <main className="container mx-auto px-4 pt-32">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
          className="text-center"
        >
          <motion.img
            src="https://raw.githubusercontent.com/HelixY2J/firefly/main/res/logo.svg"
            alt="Firefly Logo"
            className="w-24 h-24 mx-auto mb-8"
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
