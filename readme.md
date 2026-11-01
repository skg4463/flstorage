# flstorage

연합학습(Federated Learning) **기여도 평가 기반 블록체인**의 **사이드체인(Layer 2, 스토리지 체인)** 입니다.
Cosmos SDK + Ignite CLI로 만들었으며, L-node가 제출한 학습 결과(가중치 파일)와 CL-node의 글로벌 모델 파일에 대해 **원본 해시, 태그, 샤드 해시 목록**을 온체인에 기록합니다.

> 서비스 배포용이 아니라 **논문 실험 환경 구축**을 목표로 하는 연구용 프로젝트입니다.

## 프로젝트 구성

이 프로젝트는 역할별로 3개의 레포지토리로 나뉘어 있습니다.

| 레포지토리 | 역할 | 주요 모듈/구성 |
|---|---|---|
| [flmainchain](https://github.com/skg4463/flmainchain) | 메인체인. 라운드·위원회·기여도(ATT)·글로벌 모델 해시 관리 | `x/fedlearning` |
| **[flstorage](https://github.com/skg4463/flstorage)** (이 레포) | 사이드체인. 학습 결과 원본(가중치 파일)의 메타데이터·샤드 해시 저장 | `x/fedstoraging` |
| [fls-client](https://github.com/skg4463/fls-client) | 오프체인 클라이언트(업로더/다운로더), 자동 시뮬레이션·성능 측정 스크립트 | `uploader/`, `downloader/`, `simulate_auto.sh` |

전체 아키텍처 다이어그램은 [flmainchain의 diagram](https://github.com/skg4463/flmainchain/blob/main/diagram.png)을 참고하세요.

## 이 체인이 하는 일

수십 MB 크기의 학습 결과 파일을 메인체인에 직접 올리지 않기 위해, 파일 본문은 체인 밖에 두고 **무결성 검증과 복원에 필요한 메타데이터만** 이 체인에 저장합니다.

```
[L-node / CL-node]                     [flstorage]                    [flmainchain]
 weight.bin
   │ fls-client uploader
   ├─ SHA-256 → originalHash
   ├─ Reed-Solomon(10+5) → 15개 샤드 ──► store-file(tag, originalHash, shardHashes)
   │                                                 
   └──────────── originalHash, tag ─────────────────────────────────► submit-weight

[C-node]
   fls-client downloader(originalHash) ─► show-stored-file → 샤드 해시 목록
   └─ 샤드 10개 이상으로 복원 → SHA-256 비교로 무결성 검증
```

- **태그 형식**: `ROUND-USER-CHAINID` (예: `1-node0-flstorage`), 글로벌 모델은 `ROUND-global-flstorage`
- 같은 `originalHash`로 메인체인의 `ModelSubmission`과 이 체인의 `StoredFile`이 연결됩니다.

## `x/fedstoraging` 모듈

### 트랜잭션

| 명령 | 설명 |
|---|---|
| `store-file [tag] [original-hash] [shard-hashes...]` | 파일 메타데이터 저장. 보통 fls-client `uploader`가 호출 |
| `request-data-access [original-hash]` | C-node의 원본 접근 권한 요청 (현재는 기록 없이 성공만 반환, 아래 [구현 현황](#구현-현황--한계) 참고) |

### 쿼리 / 상태

| 상태 | 키 | 내용 | 쿼리 |
|---|---|---|---|
| `StoredFile` | original hash | `original_hash`, `tag`, `shard_hashes`, `creator` | `show-stored-file <hash>`, `list-stored-file` |
| `DataAccessPermission` | permission id | `granted` 여부 | `show-data-access-permission <id>`, `list-data-access-permission` |

(`show-*`는 `get-*`의 별칭입니다.)

### IBC

- 포트 `fedstoraging`
- 메인체인(`fedlearning` 포트)과 채널을 맺을 때 채널 버전 `fedlearning-1`을 허용하도록 `OnChanOpenTry`를 수정해 두었습니다.
- 채널 생성 방법은 [flmainchain README의 릴레이어 설정](https://github.com/skg4463/flmainchain#2-ibc-릴레이어hermes-설정)을 참고하세요.

`x/flstorage`는 Ignite가 기본으로 생성한 모듈(params만 있음)로, 프로토콜 로직은 모두 `x/fedstoraging`에 있습니다.

## 개발 환경

| 항목 | 버전 |
|---|---|
| Ignite CLI | v29.3.1-dev |
| Cosmos SDK | v0.53.3 |
| ibc-go | v10 |
| Go | 1.24.7 |

- 바이너리: `flstoraged`, 체인 ID: `flstorage`
- 주소 prefix: `fls` (메인체인과 동일)

## 실행

### 1. 체인 실행

```bash
ignite chain serve
```

이 체인은 Ignite 기본 포트를 사용하므로 `flstoraged` 명령에 `--node`를 붙이지 않아도 됩니다. (메인체인은 26658 등 별도 포트 사용)

| 용도 | flstorage | (참고) flmainchain |
|---|---|---|
| CometBFT RPC | **26657** | 26658 |
| REST API | **1317** | 1318 |
| gRPC | **9090** | 9091 |
| Faucet | **4500** | 4501 |

`config.yml`의 `alice`, `bob`, `node0`~`node9` 계정은 메인체인과 **같은 니모닉**을 사용하므로, 같은 키 이름은 두 체인에서 같은 `fls1...` 주소를 가집니다.

```bash
flstoraged keys list
```

### 2. 파일 업로드 → 조회 → 복원 확인

[fls-client](https://github.com/skg4463/fls-client)를 이 레포와 같은 상위 디렉터리에 받아 두었다고 가정합니다.

**① 업로드** (`fls-client/uploader`에서 실행) — 결과에 `Original file hash`와 txhash가 출력됩니다.

```bash
cd ../fls-client/uploader
go run main.go ../model_round_1_client_1.bin "1-alice-flstorage" alice
```

**② 온체인 메타데이터 조회**

```bash
flstoraged query fedstoraging show-stored-file <originalHash>
```

```bash
flstoraged query tx <txHash>
```

```bash
flstoraged query block --type=height <height>
```

**③ 복원** (`fls-client/downloader`에서 실행)

```bash
cd ../downloader
go run main.go <originalHash> ../restored_round_weight.bin
```

**④ 원본과 복원 파일 비교**

```bash
cd ..
sha256sum model_round_1_client_1.bin restored_round_weight.bin
```

## 디렉터리 구조

```
flstorage/
├── app/                        # 앱 구성, IBC 라우터 등록 (app/ibc.go)
├── cmd/flstoraged/             # 노드 바이너리
├── proto/flstorage/fedstoraging/v1/   # 메시지·상태 정의
├── x/fedstoraging/
│   ├── keeper/msg_server_store_file.go          # 메타데이터 저장
│   ├── keeper/msg_server_request_data_access.go # 접근 권한 요청 (TODO: ICQ)
│   └── module/module_ibc.go                     # IBC 채널 핸드셰이크
├── x/flstorage/                # Ignite 기본 모듈
└── config.yml                  # Ignite 체인 설정 (계정, validator)
```

## 참고

- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)


## Get started

```
ignite chain serve

계정 및 키 확인
flstoraged keys list

업로드 확인된 alice 키 삽입 [fls-client - uploader]
go run main.go ../model_round_1_client_1.bin "1-[aliceAddr]-flstorage" [aliceAddr]

결과의 originalHash를 통해 query 
flstoraged query storage show-stored-file [originalHash]

txhash를 통해 블록 생성 확인 
flstoraged query tx [txHash]

block height를 통한 블록 확인
flstoraged query block --type=height 30248

downloader를 통해 originalHash를 인자로 복원 [fls-client - downloader]
go run main.go [originalHash] ../restored_round_weight.bin

checksum 비교를 통해 복원 확인
sha256sum model_round_1_client_1.bin restored_round_weight.bin
```
