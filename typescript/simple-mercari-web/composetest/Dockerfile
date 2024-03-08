# syntax=docker/dockerfile:1
FROM python:3.7-alpine
WORKDIR /code
ENV FLASK_APP=app.py
ENV FLASK_RUN_HOST=0.0.0.0
RUN apk add --no-cache gcc musl-dev linux-headers
COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt
EXPOSE 5000
COPY . .
CMD ["flask", "run"]

# Python 3.7 イメージを使って当イメージを構築する。
# 作業用ディレクトリを /code に指定する。
# flaskコマンドにより用いられる環境変数を設定する。
# gcc や依存パッケージをインストールする。
# requirements.txtをコピーして Python 依存パッケージをインストールする。
# イメージにメタデータを追加して、コンテナーがポート 5000 をリッスンするように記述する。
# このプロジェクトのカレントディレクトリ.を、イメージ内のワークディレクトリ.にコピーする。
# コンテナーに対するデフォルトのコマンドを flask run にする。
