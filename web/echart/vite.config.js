import vue from '@vitejs/plugin-vue'

import { resolve } from 'path'

const pathResolve = (dir) => {
  return resolve(__dirname, ".", dir)
}

const alias = {
  '@': pathResolve("src")
}

// https://vitejs.dev/config/
export default ({ command }) => {
 
  return {
    base: './',
    resolve: {
      alias
    },
    server: {
      port: 3002,
      host: '0.0.0.0',
      open: true,
      // proxy: { // 代理配置   
      //   [process.env.VITE_BASE_URL]: {
      //     target:'http://192.168.0.181:9092/',
      //     changeOrigin: true,
      //   }
      // },
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks: {
            
          }
        }
      }
    },
    plugins: [
      vue(),
   
    ]
  };
}
