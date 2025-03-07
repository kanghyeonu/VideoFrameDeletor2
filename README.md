# VideoFrameDeletor2
2025 / 2 / 4 ~ 2025 / 3 / 6

**VideoFrameDeletor repository 코드 리팩토링** 프로젝트



H.264/AVC 비디오 압축 표준으로 PIR로 인코딩된 프레임을 원하는 범위를 바이트레벨에서 삭제하는 프로그램 
(논문: Frame Sceduling Approach for Real-Time Streaming(ICTC 2024))

---

### Requirement

- go 1.23.5
- 비디오 변환 툴(i.e FFmpeg) -> mp4, avi -> h264

---
### Class Diagram

![Image](https://github.com/user-attachments/assets/6425764a-1cbb-44af-aabb-73d76a60be06)
---

### Usage

- 코드 다운로드

  ```
  git clone https://github.com/kanghyeonu/VideoFrameDeletor2.git

- original videos 디렉토리에 수정할 원본 비디오 파일(.h264) 저장 

- 파라미터

  1. **filename** : Input file name with .h264 extension
  2. **bytesToRemove** : 0 to 100 if ratio is true, 0 to n if ratio is false
  3. **start offset** : 0 to 100 offset starting position for deletion in each Nalu
  4. **ratio** : Ratio for processing (true: 1/false: 0) 
  5. **reverse** : Reverse the operation (true: 1/false: 0)
  6. **increment** : Increment value for start offset

- 예시

  ```
  ~/VideoFrameDeletor2> go run main.go BBB_PIR.h264 5 5 1 0 5
  ```

  1. filename: BBB_PIR.h264
  2. bytesToRemove: 5 <- 삭제할 바이트 수 5 bytes or 5%
  3. start offset: 5 <- 삭제를 시작할 위치 nalu 크기의 5%되는 byte index위치
  4. ratio: 1 <- **참**이면 비율기반으로 바이트 삭제, **거짓**이면 상수값으로 삭제
  5. reverse: 0 <- **참**이면 정방향 좌 -> 우, **거짓**이면 역방향 우 -> 좌
  6. incremet: 5 <- start offset의 증가량 5%, 10%, 15% ...에서 삭제 시작위치

- 결과 파일

  **예시** **결과 파일의 경우 modified vidoes/5_5_true_false_5/ 의 경로**로 저장

  -> offset5.h264, offset10.h264, offset15.h264..., offset95.h264 파일 생성

---

### 동작 알고리즘

- h264 표준에서 각 Nalu의 시작 패턴은 **3byte 패턴(0x00001)**과 **4byte 패턴(0x000001)**으로 구성
- 해당 시작 패턴을 찾아서 각 Nalu의 크기를 기반으로 일정 수의 바이트 삭제 (비디오 메타데이터 Nalu는 제외)

---
### 결과물
<p align="center">
  <img src="https://github.com/user-attachments/assets/e48a3fff-27f9-4817-96ef-74f70928510a" align="center" width="50%">
  <img src="https://github.com/user-attachments/assets/9cb745bf-23d8-479a-a733-5ce04e31b11d" align="center" width="50%">  
  <figcation align = "cennter">
</p>
- 손상된 프레임을 디코딩 후 프레임 추출 결과(위: 원본 프레임, 아래: 200byte가 삭제된 프레임)
  



