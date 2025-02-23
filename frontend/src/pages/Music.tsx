
import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { Music as MusicIcon, Play, Pause, SkipBack, SkipForward, Volume2, AlertCircle } from 'lucide-react';
import { Navbar } from '../components/Navbar';

// Sample hardcoded music data
const SAMPLE_SONGS = [
  { id: 1, title: "Midnight City", artist: "M83", duration: "4:03" },
  { id: 2, title: "Shadows", artist: "Roosevelt", duration: "3:45" },
  { id: 3, title: "Instant Crush", artist: "Daft Punk", duration: "5:37" },
  { id: 4, title: "Gold", artist: "Chet Faker", duration: "4:48" },
  { id: 5, title: "Electric Feel", artist: "MGMT", duration: "3:49" },
];

const MusicPage = () => {
  const [songs, setSongs] = useState(SAMPLE_SONGS);
  const [currentSong, setCurrentSong] = useState(SAMPLE_SONGS[0]);
  const [isPlaying, setIsPlaying] = useState(false);
  const [wsConnected, setWsConnected] = useState(false);

  useEffect(() => {
    // Create WebSocket connection
    const ws = new WebSocket('ws://localhost:8081/ws');

    ws.onopen = () => {
      setWsConnected(true);
      console.log('Connected to music WebSocket');
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      
      // Handle different message types
      switch (data.type) {
        case "master_stats":
          console.log("Memory Stats:", data.memory);
          console.log("MP3 Files:", data.mp3_files);
          // Update songs list with received MP3 files
          const newSongs = data.mp3_files.map((filename: string, index: number) => ({
            id: index + 1,
            title: filename.replace('.mp3', ''),
            artist: 'Unknown', // You can modify this as needed
            duration: "0:00"   // You can modify this as needed
          }));
          setSongs(newSongs);
          break;
          
        default:
          // Handle other message types if needed
          if (data.songs) {
            setSongs(data.songs);
          }
          break;
      }
    };

    ws.onclose = () => {
      setWsConnected(false);
      console.log('Disconnected from music WebSocket');
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <div className="min-h-screen">
      <Navbar />
      
      <main className="container mx-auto px-4 pt-32 pb-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Playlist Section */}
          <div className="glass-panel p-6">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold flex items-center gap-2">
                <MusicIcon />
                Available Songs
              </h2>
              {!wsConnected && (
                <div className="flex items-center gap-2 text-yellow-400 text-sm">
                  <AlertCircle size={16} />
                  <span>Demo Mode</span>
                </div>
              )}
            </div>
            
            <div className="space-y-2">
              {songs.map((song) => (
                <motion.div
                  key={song.id}
                  className={`p-3 rounded-lg cursor-pointer transition-colors ${
                    currentSong.id === song.id
                      ? 'bg-glow/20 text-white'
                      : 'hover:bg-white/5 text-white/70'
                  }`}
                  onClick={() => setCurrentSong(song)}
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                >
                  <div className="flex justify-between items-center">
                    <div>
                      <h3 className="font-medium">{song.title}</h3>
                      <p className="text-sm text-white/50">{song.artist}</p>
                    </div>
                    <span className="text-sm text-white/30">{song.duration}</span>
                  </div>
                </motion.div>
              ))}
            </div>
          </div>
          
          {/* Music Player Section */}
          <div className="lg:col-span-2">
            <div className="glass-panel p-8 h-full flex flex-col">
              {/* Album Art / Vinyl */}
              <div className="flex-1 flex items-center justify-center mb-8">
                <motion.div
                  className="w-64 h-64 rounded-full bg-glow/20 flex items-center justify-center relative"
                  animate={{
                    rotate: isPlaying ? 360 : 0,
                  }}
                  transition={{
                    duration: 3,
                    repeat: Infinity,
                    ease: "linear",
                  }}
                >
                  <div className="absolute inset-4 rounded-full bg-background" />
                  <div className="absolute inset-[18px] rounded-full bg-white/5 backdrop-blur" />
                  <div className="absolute inset-0 rounded-full border-2 border-glow/30" />
                  <div className="absolute inset-[70px] rounded-full bg-glow/5">
                    <div className="absolute inset-0 flex items-center justify-center">
                      <MusicIcon size={48} />
                    </div>
                  </div>
                </motion.div>
              </div>
              
              {/* Song Info */}
              <div className="text-center mb-8">
                <h2 className="text-2xl font-semibold mb-2">{currentSong.title}</h2>
                <p className="text-white/60">{currentSong.artist}</p>
              </div>
              
              {/* Progress Bar */}
              <div className="w-full bg-white/5 rounded-full h-1 mb-8">
                <div className="bg-glow h-full rounded-full" style={{ width: '45%' }} />
              </div>
              
              {/* Controls */}
              <div className="flex items-center justify-center space-x-8">
                <button className="p-2 text-white/60 hover:text-white transition-colors">
                  <SkipBack size={24} />
                </button>
                <button 
                  className="p-4 bg-glow rounded-full text-background hover:bg-glow/90 transition-colors"
                  onClick={() => setIsPlaying(!isPlaying)}
                >
                  {isPlaying ? <Pause size={24} /> : <Play size={24} />}
                </button>
                <button className="p-2 text-white/60 hover:text-white transition-colors">
                  <SkipForward size={24} />
                </button>
              </div>
              
              {/* Volume */}
              <div className="flex items-center justify-end mt-6">
                <Volume2 size={18} />
                <div className="w-32 bg-white/5 rounded-full h-1">
                  <div className="bg-white/20 h-full rounded-full" style={{ width: '75%' }} />
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default MusicPage;
