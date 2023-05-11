#!/usr/bin/python3
# pip3 install bilibili-api-python
# https://www.bilibili.com/video/BV1NP411o7MB
import asyncio
from bilibili_api import video, Credential, HEADERS
import httpx
import os
import argparse

# FFMPEG 路径，查看：http://ffmpeg.org/
FFMPEG_PATH = "ffmpeg"


async def download_url(url: str, out: str, info: str):
    # 下载函数
    async with httpx.AsyncClient(headers=HEADERS) as sess:
        resp = await sess.get(url)
        length = resp.headers.get('content-length')
        with open(out, 'wb') as f:
            process = 0
            for chunk in resp.iter_bytes(1024):
                if not chunk:
                    break

                process += len(chunk)
                # print(f'下载 {info} {process} / {length}')
                f.write(chunk)


async def main(args: argparse.Namespace):
    full_name_path = f'{args.save_path}'

    # 实例化 Credential 类
    credential = Credential(sessdata=args.sessdata,
                            bili_jct=args.bili_jct, buvid3=args.buvid3)
    # 实例化 Video 类
    v = video.Video(bvid=args.bvid, credential=credential)
    if os.path.exists(f'{full_name_path}.mp4'):
        return

    # 创建文件夹 dir(full_path)
    dst_dir = os.path.dirname(full_name_path)
    if not os.path.exists(dst_dir):
        os.makedirs(dst_dir, exist_ok=True)

    # 获取视频下载链接
    download_url_data = await v.get_download_url(cid=args.cid)

    # 解析视频下载信息
    detecter = video.VideoDownloadURLDataDetecter(data=download_url_data)
    streams = detecter.detect_best_streams()
    if len(streams)!=2:
        return

    if streams[0] is None or streams[1] is None:
        return

    # 有 MP4 流 / FLV 流两种可能
    if detecter.check_flv_stream() == True:
        # FLV 流下载
        await download_url(streams[0].url, f"{full_name_path}_temp.flv", "FLV 音视频流")
        # 转换文件格式
        os.system(
            f'{FFMPEG_PATH} -loglevel quiet -i {full_name_path}_temp.flv {full_name_path}.mp4')
        # 删除临时文件
        os.remove(f"{full_name_path}_temp.flv")
    else:
        # MP4 流下载
        await download_url(streams[0].url, f"{full_name_path}_video_temp.m4s", "视频流")
        await download_url(streams[1].url, f"{full_name_path}_audio_temp.m4s", "视频流")
        # 混流
        os.system(f'{FFMPEG_PATH} -loglevel quiet -i {full_name_path}_video_temp.m4s -i {full_name_path}_audio_temp.m4s -vcodec copy -acodec copy {full_name_path}.mp4')
        # 删除临时文件
        os.remove(f"{full_name_path}_video_temp.m4s")
        os.remove(f"{full_name_path}_audio_temp.m4s")

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--bvid', required=True)
    parser.add_argument('--cid', required=True)
    parser.add_argument('--buvid3', required=True)
    parser.add_argument('--bili_jct', required=True)
    parser.add_argument('--sessdata', required=True)
    parser.add_argument('--save_path', default="/opt/downloads/bilibili")
    args = parser.parse_args()

    # print(args.save_path)
    # 主入口
    asyncio.get_event_loop().run_until_complete(main(args))
