# HaberlerPlus - Gelişmiş CLI Haber Bülteni

HaberlerPlus, çeşitli haber kaynaklarından haber başlıklarını ve URL'leri çekmek için kullanılan gelişmiş bir komut satırı aracıdır. Orijinal CLIHaberBulteni projesinin genişletilmiş versiyonudur.

## Desteklenen Haber Kaynakları

1. **GZT.com** - Orijinal haber kaynağı
2. **Hurriyet.com.tr** - Hürriyet gazetesi web sitesi
3. **Sozcu.com.tr** - Sözcü gazetesi web sitesi

## Kurulum

### Kaynak Koddan Derleme

1. Projeyi klonlayın:
   ```bash
   git clone https://github.com/furkandogmus/HaberlerPlus.git
   cd HaberlerPlus
   ```

2. Bağımlılıkları yükleyin:
   ```bash
   go mod tidy
   ```

3. Projeyi derleyin:
   ```bash
   go build -o haberlerplus ./cmd/haberlerplus
   ```

4. Binary dosyasını sistem yoluna kopyalayın:
   ```bash
   sudo cp ./haberlerplus /usr/local/bin/haberlerplus
   ```

## Kullanım

1. Programı çalıştırın:
   ```bash
   haberlerplus
   ```

2. Haber kaynağını seçin (1-3 arası bir sayı girin).
3. Seçtiğiniz haber kaynağı için kategori seçin.
4. Haberleriniz gösterilecektir!

## Özellikler

- Birden fazla haber kaynağı desteği
- Kategori bazlı haber görüntüleme
- Renkli terminal çıktısı
- Kullanıcı dostu arayüz
- Reklamsız ve hızlı haber erişimi

## Komut Satırı Seçenekleri

- `-h`: Yardım bilgisini gösterir
- `-v`: Versiyon bilgisini gösterir

## Katkıda Bulunma

Her türlü katkıya açığız! Yeni özellikler eklemek, hata düzeltmek veya mevcut kodu geliştirmek isterseniz, lütfen katkıda bulunun.

### Yeni Haber Kaynağı Ekleme

Yeni bir haber kaynağı eklemek için:

1. `pkg/sources/sources.go` dosyasında `NewsSource` arayüzünü uygulayan yeni bir yapı oluşturun.
2. `GetAllSources()` fonksiyonuna yeni kaynağınızı ekleyin.

## Lisans

Bu proje MIT lisansı altında lisanslanmıştır.

## İletişim

Furkan Doğmuş - furkandogmus9183@gmail.com 