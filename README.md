# NaluHandler

H.264/AVC 비디오 압축 표준으로 인코딩된 프레임을 삭제하는 프로그램 
(논문: ICTC 2024 Frame Sceduling Approach for Real-Time Streaming)

---

### 동작 방법

- FFmpeg와 같은 툴로 .mp4 파일을 .h264로 변환
- 파일 내 main문 시작의 for문의 start_offset 변수 설정으로 시작 위치 및 삭제 비율 설정 후 실행
  - **RatioForDeleting:** 삭제 시킬 비율 or 고정 바이트
  - **offset:** 인코딩된 각 프레임의 삭제 시작 위치 설정
  - **ratio**: 프레임의 일정 비율만큼 삭제 if Ture, 프레임을 고정 바이트 수만큼 삭제 if False
  - **reverse**: 프레임 삭제 시작 위치의 기준을 Head 또는 Tail로 설정

---

### 동작 알고리즘

- KMP 알고리즘을 통해서 NAL Unit 시작 패턴 검색
- 시작 패턴은 3byte 패턴(00001)과 4byte 패턴(000001)으로 구성
- go Rutine을 이용해 읽어들인 파일의 3, 4바이트 패턴의 위치를 검색
- 찾아낸 하나의 Nalu를 처리하여 반환
---
### 결과물
<p align="center">
  <img src="https://github.com/user-attachments/assets/e48a3fff-27f9-4817-96ef-74f70928510a" align="center" width="50%">
  <img src="https://github.com/user-attachments/assets/9cb745bf-23d8-479a-a733-5ce04e31b11d" align="center" width="50%">  
  <figcation align = "cennter">
</p>
- 손상된 프레임을 디코딩 후 프레임 추출 결과(위: 원본 프레임, 아래: 200byte가 삭제된 프레임)
  
## 리팩토리 중
