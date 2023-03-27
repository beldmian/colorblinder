import Head from 'next/head'
import {Button, Typography, Card, Col, Row,
  Space, Input, Spin, Switch, Slider, InputNumber, ConfigProvider, theme} from 'antd'
import { RgbaColorPicker } from 'react-colorful'
import dynamic from 'next/dynamic'
import { useState } from 'react'
import Filter from '@/components/Filter';

const DynamicPlayer = dynamic(() => import('../components/Player'), {
  ssr: false,
})

export default function Home() {
  let componentToHex = (c) => {
    var hex = c.toString(16);
    return hex.length == 1 ? "0" + hex : hex;
  }
  let rgbToHex = (r, g, b) => {
    return "#" + componentToHex(r) + componentToHex(g) + componentToHex(b);
  }
  let [url, setUrl] = useState("")
  let [endUrl, setEndUrl] = useState("")
  let [loading, setLoading] = useState(false)
  let [color, setColor] = useState({
      r: 255,
      g: 255,
      b: 255,
      a: 0
  })
  let [contrast, setContrast] = useState(1)
  let [saturate, setSaturate] = useState(1)
  let [isPhotosensitive, setIsPhotosensitive] = useState(false)
  let fetchData = async () => {
    console.log(color)
    setLoading(true)
    if (isPhotosensitive) {
      const req = await fetch('http://127.0.0.1:8080/start_stream', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          stream_url: url
        })
      })
      const res = await req.json()
      while (true) {
        const req2= await fetch("http://127.0.0.1:8080"+res.new_url)
        if (req2.status === 200) {
          setEndUrl("http://127.0.0.1:8080"+res.new_url)
          setLoading(false)
          break
        }
        await new Promise(r => setTimeout(r, 500))
      }
      return
    }
    setEndUrl(url)
    setLoading(false)
  }
  return (
    <>
      {/* {color.rgb != undefined && ( */}
        <div style={{opacity: 0}}>
          <Filter color={rgbToHex(color.r, color.g, color.b)} color_opacity={color.a} saturation={saturate} contrast={contrast} />
        </div>
      {/* )} */}
      <Head>
        <title>Colorblinder</title>
      </Head>
      <ConfigProvider
        theme={{
          // token: {
          //   colorPrimary: '#00b96b',
          // },
          algorithm: theme.darkAlgorithm,
        }}
      >
      <Row gutter={16}>
        <Col span={12}>
          <Card title="Настройки" size="default">
            <Space direction='vertical'>
              <Typography.Title level={3}>Настройки фильтра</Typography.Title>
              <Space>
                <Typography>Контрастность</Typography>
                <Slider
                style={{width: "200px"}}
                  min={0}
                  max={5}
                  onChange={(value) => setContrast(value)}
                  value={contrast}
                  step={0.01}
                />
              </Space>
              <Space>
                <Typography>Насыщенность</Typography>
                <Slider
                style={{width: "200px"}}
                  min={0}
                  max={5}
                  onChange={(value) => setSaturate(value)}
                  value={saturate}
                  step={0.01}
                />
              </Space>
              <Space>
                <RgbaColorPicker color={color} onChange={setColor} />
                <div style={{
                  width: "200px",
                  height: "200px",
                  backgroundImage: `url("/img/colorwheel.png")`,
                  backgroundSize: "cover",
                  filter: "url(#video_filter)",
                }}></div>
              </Space>
              <Space>
                <Switch onChange={(checked) => setIsPhotosensitive(checked)}/>
                <Typography>Фильтрация сцен, которые могут вызвать эпилептический припадок</Typography>
              </Space>
              <Typography.Title level={3}>Настройки потока ввода</Typography.Title>
              <Input placeholder='URL-адрес mpd-потока' onChange={(e) => setUrl(e.target.value)}></Input>
              <Button type="default" onClick={(e) => {
                e.preventDefault()
                fetchData()
              }}>
                {loading ? <Spin size="small" /> : 'Загрузить'}
              </Button>
            </Space>
          </Card>
        </Col>
        <Col span={12}>
          <Card title="Плеер" size="default">
              <DynamicPlayer url={endUrl} />
          </Card>
        </Col>
      </Row>
      </ConfigProvider>
    </>
  )
}
