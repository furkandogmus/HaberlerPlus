# HaberlerPlus - Gelişmiş CLI Haber Bülteni

HaberlerPlus, çeşitli haber kaynaklarından haber başlıklarını ve URL'leri çekmek için kullanılan gelişmiş bir komut satırı aracıdır. Orijinal CLIHaberBulteni projesinin genişletilmiş versiyonudur.

## Desteklenen Haber Kaynakları

### HTML Tabanlı Kaynaklar
1. **GZT.com** - Orijinal haber kaynağı
2. **Hurriyet.com.tr** - Hürriyet gazetesi web sitesi
3. **Sozcu.com.tr** - Sözcü gazetesi web sitesi
4. **Milliyet.com.tr** - Milliyet gazetesi web sitesi
5. **Haberler.com** - Haberler.com haber portalı

### RSS Tabanlı Kaynaklar
6. **CNN Türk** - CNN Türk RSS beslemeleri (GÜNDEM, DÜNYA, EKONOMİ, SPOR, SAĞLIK, TEKNOLOJİ)
7. **NTV** - NTV RSS beslemeleri (SON DAKİKA, GÜNDEM, DÜNYA, EKONOMİ, SPOR, SAĞLIK, TEKNOLOJİ)
8. **Habertürk** - Habertürk RSS beslemeleri (GÜNDEM, DÜNYA, EKONOMİ, SPOR, SAĞLIK, TEKNOLOJİ)

## Kurulum

### Doğrudan Kurulum (Go Kullanıcıları İçin)

Go yüklü ise, aşağıdaki komutla doğrudan kurulum yapabilirsiniz:

```bash
go install github.com/furkandogmus/HaberlerPlus/cmd/news@latest
```

Bu komut, uygulamayı `$GOPATH/bin` dizinine kuracaktır. Bu dizinin PATH'inizde olduğundan emin olun.

### Önceden Derlenmiş Binary İndirme

Önceden derlenmiş binary dosyalarını [Releases](https://github.com/furkandogmus/HaberlerPlus/releases) sayfasından indirebilirsiniz. İndirdiğiniz binary dosyasını çalıştırılabilir yapın ve PATH'inizde olan bir dizine kopyalayın:

```bash
chmod +x news
sudo mv news /usr/local/bin/
```

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
   go build -o bin/news ./cmd/news
   ```

4. Binary dosyasını sistem yoluna kopyalayın (opsiyonel):
   ```bash
   sudo cp ./bin/news /usr/local/bin/news
   ```

## Kullanım

1. Programı çalıştırın:
   ```bash
   news
   ```

2. Haber kaynağını seçin (1-8 arası bir sayı girin).
3. Seçtiğiniz haber kaynağı için kategori seçin.
4. Haberleriniz gösterilecektir!

## Komut Satırı Seçenekleri

- `-h`: Yardım bilgisini gösterir
- `-v`: Versiyon bilgisini gösterir
- `-d`: Debug modunda çalıştırır

## Desteklenen Kategoriler

Her haber kaynağı için desteklenen kategoriler:

### HTML Tabanlı Kaynaklar
- **GZT.com**: GÜNDEM, DÜNYA, EKONOMİ, SPOR, TEKNOLOJİ
- **Hurriyet.com.tr**: GÜNDEM, DÜNYA, EKONOMİ, SPOR, TEKNOLOJİ
- **Sozcu.com.tr**: GÜNDEM, DÜNYA, EKONOMİ, SPOR, SAĞLIK, TEKNOLOJİ
- **Milliyet.com.tr**: GÜNDEM, DÜNYA, EKONOMİ, SPOR, TEKNOLOJİ
- **Haberler.com**: GÜNDEM, DÜNYA, EKONOMİ, SPOR, SAĞLIK, TEKNOLOJİ

### RSS Tabanlı Kaynaklar
- **CNN Türk**: GÜNDEM, DÜNYA, EKONOMİ, SPOR, SAĞLIK, TEKNOLOJİ
- **NTV**: SON DAKİKA, GÜNDEM, DÜNYA, EKONOMİ, SPOR, SAĞLIK, TEKNOLOJİ
- **Habertürk**: GÜNDEM, DÜNYA, EKONOMİ, SPOR, SAĞLIK, TEKNOLOJİ

## Geliştirici Bilgileri

### Kütüphane Olarak Kullanım

HaberlerPlus'ı kendi Go projelerinizde kütüphane olarak kullanabilirsiniz:

```go
package main

import (
    "fmt"
    "github.com/furkandogmus/HaberlerPlus/pkg/sources"
)

func main() {
    // Tüm haber kaynaklarını al
    allSources := sources.GetAllSources()
    
    // İlk kaynağı seç
    source := allSources[0]
    
    // Kategorileri göster
    fmt.Printf("Kategoriler: %v\n", source.Categories())
    
    // İlk kategoriden haberleri getir
    news, err := source.FetchNews(0)
    if err != nil {
        panic(err)
    }
    
    // Haberleri göster
    for _, item := range news {
        fmt.Printf("%s: %s\n", item.Title, item.URL)
    }
}
```

## Katkıda Bulunma

Her türlü katkıya açığız! Yeni özellikler eklemek, hata düzeltmek veya mevcut kodu geliştirmek isterseniz, lütfen katkıda bulunun.

### Yeni Haber Kaynağı Ekleme

Yeni bir haber kaynağı eklemek için:

1. HTML tabanlı kaynaklar için:
   - `pkg/sources/impl` dizininde yeni bir kaynak dosyası oluşturun.
   - `NewsSource` arayüzünü uygulayan bir yapı oluşturun.

2. RSS tabanlı kaynaklar için:
   - `pkg/sources/impl` dizininde yeni bir RSS kaynak dosyası oluşturun.
   - `RSSSource` yapısını kullanarak yeni bir kaynak oluşturun veya kendi özel yapınızı oluşturun.

3. `pkg/sources/factory.go` dosyasına yeni bir factory fonksiyonu ekleyin.
4. `pkg/sources/sources.go` dosyasındaki `GetAllSources()` fonksiyonuna yeni kaynağınızı ekleyin.

## Test Etme

Tüm haber kaynaklarını test etmek için:

```bash
./scripts/test_all_sources.sh
```

Bu script, tüm haber kaynaklarını ve kategorilerini test eder ve sonuçları gösterir.

## Lisans

Bu proje MIT lisansı altında lisanslanmıştır.

## İletişim

Furkan Doğmuş - furkandogmus9183@gmail.com 