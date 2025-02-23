
import { Link, useLocation } from 'react-router-dom';
import { Music, Network } from 'lucide-react';

export const Navbar = () => {
  const location = useLocation();
  
  return (
    <nav className="fixed top-0 left-0 right-0 z-50">
      <div className="glass-panel mx-4 my-4 px-6 py-4">
        <div className="flex items-center justify-between">
          <Link to="/" className="flex items-center space-x-2">
            <img src="https://raw.githubusercontent.com/HelixY2J/firefly/main/res/logo.svg" 
                 alt="Firefly" 
                 className="h-8 w-8" />
            <span className="font-semibold text-lg">Firefly</span>
          </Link>
          
          <div className="flex items-center space-x-6">
            <Link 
              to="/music" 
              className={`nav-link flex items-center space-x-2 ${
                location.pathname === '/music' ? 'text-glow' : ''
              }`}
            >
              <Music size={18} />
              <span>Music</span>
            </Link>
            <Link 
              to="/nodes" 
              className={`nav-link flex items-center space-x-2 ${
                location.pathname === '/nodes' ? 'text-glow' : ''
              }`}
            >
              <Network size={18} />
              <span>Nodes</span>
            </Link>
          </div>
        </div>
      </div>
    </nav>
  );
};
