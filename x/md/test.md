制作过程
1.准备环境
2. 准备打包工具
3. 准备文件
4. 打包
安装过程
参考：

1. https://www.pugetsystems.com/labs/hpc/ubuntu-22-04-server-autoinstall-iso/

2. https://github.com/covertsh/ubuntu-autoinstall-generator/tree/main

3. https://canonical-subiquity.readthedocs-hosted.com/en/latest/reference/autoinstall-reference.html#late-commands

4. https://help.ubuntu.com/community/LiveCDCustomization



制作过程
1.准备环境
一台ubuntu 服务器（安装对应的软件包：7z， xorriso）
官方ISO镜像（使用ubuntu-22.04.1-live-server-amd64.iso， 固定这个镜像，MD5： e8d2a77c51b599c10651608a5d8c286f）
2. 准备打包工具
2.0)  mkdir -p iso ; cd iso; wget https://xxxxxxxxxxxxxx/ubuntu-22.04.1-live-server-amd64.iso   (官方镜像)

2.1） git clone https://github.com/covertsh/ubuntu-autoinstall-generator.git

2.2）修改代码：

diff --git a/ubuntu-autoinstall-generator.sh b/ubuntu-autoinstall-generator.sh
index 5229d83..d9dbab3 100644
--- a/ubuntu-autoinstall-generator.sh
+++ b/ubuntu-autoinstall-generator.sh
@@ -4,7 +4,7 @@ set -Eeuo pipefail
 function cleanup() {
         trap - SIGINT SIGTERM ERR EXIT
         if [ -n "${tmpdir+x}" ]; then
-                rm -rf "$tmpdir"
+                #rm -rf "$tmpdir"
                 log " Deleted temporary working directory $tmpdir"
         fi
 }
@@ -141,7 +141,9 @@ ubuntu_gpg_key_id="843938DF228D22F7B3742BC0D94AA3F0EFE21092"
 
 parse_params "$@"
 
-tmpdir=$(mktemp -d)
+#tmpdir=$(mktemp -d)
+mkdir -p tmpdir
+tmpdir="./tmpdir"
 
 if [[ ! "$tmpdir" || ! -d "$tmpdir" ]]; then
         die " Could not create temporary working directory."
@@ -228,7 +230,7 @@ if [ ${use_hwe_kernel} -eq 1 ]; then
 fi
 
 log "Adding autoinstall parameter to kernel command line..."
-sed -i -e 's/---/ autoinstall  ---/g' "$tmpdir/isolinux/txt.cfg"
+#sed -i -e 's/---/ autoinstall  ---/g' "$tmpdir/isolinux/txt.cfg"
 sed -i -e 's/---/ autoinstall  ---/g' "$tmpdir/boot/grub/grub.cfg"
 sed -i -e 's/---/ autoinstall  ---/g' "$tmpdir/boot/grub/loopback.cfg"
 log " Added parameter to UEFI and BIOS kernel command lines."
@@ -242,7 +244,7 @@ if [ ${all_in_one} -eq 1 ]; then
         else
                 touch "$tmpdir/nocloud/meta-data"
         fi
-        sed -i -e 's,---, ds=nocloud;s=/cdrom/nocloud/  ---,g' "$tmpdir/isolinux/txt.cfg"
+        #sed -i -e 's,---, ds=nocloud;s=/cdrom/nocloud/  ---,g' "$tmpdir/isolinux/txt.cfg"
         sed -i -e 's,---, ds=nocloud\\\;s=/cdrom/nocloud/  ---,g' "$tmpdir/boot/grub/grub.cfg"
         sed -i -e 's,---, ds=nocloud\\\;s=/cdrom/nocloud/  ---,g' "$tmpdir/boot/grub/loopback.cfg"
         log " Added data and configured kernel command line."
@@ -263,8 +265,8 @@ fi
 
 log " Repackaging extracted files into an ISO image..."
 cd "$tmpdir"
-xorriso -as mkisofs -r -V "ubuntu-autoinstall-$today" -J -b isolinux/isolinux.bin -c isolinux/boot.cat -no-emul-boot -boot-load-size 4 -isohybrid-mbr /usr/lib/ISOLINUX/isohdpfx.bin -boot-info-table -input-charset utf-8 -eltorito-alt-boot -e boot/grub/efi.img -no-emul-boot -isohybrid-gpt-basdat -o "${destinat
ion_iso}" . &>/dev/null
+#xorriso -as mkisofs -r -V "ubuntu-autoinstall-$today" -J -b isolinux/isolinux.bin -c isolinux/boot.cat -no-emul-boot -boot-load-size 4 -isohybrid-mbr /usr/lib/ISOLINUX/isohdpfx.bin -boot-info-table -input-charset utf-8 -eltorito-alt-boot -e boot/grub/efi.img -no-emul-boot -isohybrid-gpt-basdat -o "${destination_iso}" . &>/dev/null
 cd "$OLDPWD"
 log " Repackaged into ${destination_iso}"
 
-die " Completed." 0
+####die " Completed." 0


当前的目录结构如下：

.
├── ubuntu-22.04.1-live-server-amd64.iso
└── ubuntu-autoinstall-generator
    ├── LICENSE
    ├── README.md
    ├── ubuntu-autoinstall-generator.sh
    └── user-data.example
 
1 directory, 5 files


3. 准备文件
3.1） 解压ISO：7z -y x ubuntu-22.04.1-live-server-amd64.iso -osource-files 

3.2） 拷贝“[BOOT]”目录： mv ./source-files/'[BOOT]' ./BOOT

3.3)  准备userdata文件：

#cloud-config
autoinstall:
  version: 1
  interactive-sections:  # Install groups listed here will wait for user input
   - storage
   - network
  storage:  # This should set the interactive (lvm set) default
    layout:
      name: lvm
      match:
        size: largest
  locale: en_US.UTF-8
  keyboard:
    layout: us
  identity:
    hostname: admin-000
    password: $y$j9T$9VwN1/fLH7LWTC.kxT67j1$1KtR1boA8cOLerKJ1Nbx.q2Z7A0f9/Ed61kCkwD9Zc0
    username: admin
  ssh:
    allow-pw: false
    install-server: false
  updates: security
  apt:
    disable_suites: [backports,security]
    primary:
      - arches: [default]
        #uri: http://us.archive.ubuntu.com/ubuntu/
        uri: http://127.0.0.1/ubuntu/
    fallback: continue-anyway
  late-commands:
    - cat /cdrom/nocloud/dcloud.tar.gz.* | tar -zxv -C /target/opt/
#  user-data: # Commands here run during first boot (cannot be interactive)
#    runcmd:
#      - cd /opt/dcloud && bash -x install.sh
此时的目录结构如下（多了source-files 、 BOOT 、 user-data）：

.
├── BOOT
│   ├── 1-Boot-NoEmul.img
│   └── 2-Boot-NoEmul.img
├── source-files
│   ├── boot
│   ├── [BOOT]
│   ├── boot.catalog
│   ├── casper
│   ├── dists
│   ├── EFI
│   ├── install
│   ├── md5sum.txt
│   └── pool
├── ubuntu-22.04.1-live-server-amd64.iso
├── ubuntu-autoinstall-generator
│   ├── LICENSE
│   ├── README.md
│   ├── ubuntu-autoinstall-generator.sh
│   └── user-data.example
└── user-data
 
1 directory, 6 files
3.4) 生产ISO物料：

cd ubuntu-autoinstall-generator
bash ubuntu-autoinstall-generator.sh -k -a -u ../user-data -s ../ubuntu-22.04.1-live-server-amd64.iso
此时的目录结构如下（多了tmpdir）：

.
├── BOOT
│   ├── 1-Boot-NoEmul.img
│   └── 2-Boot-NoEmul.img
├── source-files
│   ├── boot
│   ├── [BOOT]
│   ├── boot.catalog
│   ├── casper
│   ├── dists
│   ├── EFI
│   ├── install
│   ├── md5sum.txt
│   └── pool
├── ubuntu-22.04.1-live-server-amd64.iso
├── ubuntu-autoinstall-generator
│   ├── 843938DF228D22F7B3742BC0D94AA3F0EFE21092.keyring
│   ├── LICENSE
│   ├── README.md
│   ├── SHA256SUMS-2024-09-20
│   ├── SHA256SUMS-2024-09-20.gpg
│   ├── tmpdir
│   ├── ubuntu-autoinstall-generator.sh
│   └── user-data.example
└── user-data
 
11 directories, 13 files
3.5) 准备dcloud.tar.gz

dcloud.tar.gz 单个文件不能大于4GB（刻录iso要求），所以先针对dcloud.tar.gz进行拆分成3GB文件：

split -b 3000M -d -a 1 dcloud.tar.gz dcloud.tar.gz.
 
# 把拆分后的dcloud.tar.gz.* 文件拷贝到 tmpdir 下面
mv dcloud.tar.gz.* ubuntu-autoinstall-generator/tmpdir/


此时ubuntu-autoinstall-generator下面的目录结构：

tree -L 3 ubuntu-autoinstall-generator/
ubuntu-autoinstall-generator/
├── 843938DF228D22F7B3742BC0D94AA3F0EFE21092.keyring
├── LICENSE
├── README.md
├── SHA256SUMS-2024-09-20
├── SHA256SUMS-2024-09-20.gpg
├── tmpdir
│   ├── boot
│   │   ├── grub
│   │   └── memtest86+.bin
│   ├── boot.catalog
│   ├── casper
│   │   ├── filesystem.manifest
│   │   ├── filesystem.size
│   │   ├── initrd
│   │   ├── install-sources.yaml
│   │   ├── ubuntu-server-minimal.manifest
│   │   ├── ubuntu-server-minimal.size
│   │   ├── ubuntu-server-minimal.squashfs
│   │   ├── ubuntu-server-minimal.squashfs.gpg
│   │   ├── ubuntu-server-minimal.ubuntu-server.installer.generic.manifest
│   │   ├── ubuntu-server-minimal.ubuntu-server.installer.generic.size
│   │   ├── ubuntu-server-minimal.ubuntu-server.installer.generic.squashfs
│   │   ├── ubuntu-server-minimal.ubuntu-server.installer.generic.squashfs.gpg
│   │   ├── ubuntu-server-minimal.ubuntu-server.installer.manifest
│   │   ├── ubuntu-server-minimal.ubuntu-server.installer.size
│   │   ├── ubuntu-server-minimal.ubuntu-server.installer.squashfs
│   │   ├── ubuntu-server-minimal.ubuntu-server.installer.squashfs.gpg
│   │   ├── ubuntu-server-minimal.ubuntu-server.manifest
│   │   ├── ubuntu-server-minimal.ubuntu-server.size
│   │   ├── ubuntu-server-minimal.ubuntu-server.squashfs
│   │   ├── ubuntu-server-minimal.ubuntu-server.squashfs.gpg
│   │   └── vmlinuz
│   ├── dists
│   │   ├── jammy
│   │   ├── stable -> jammy
│   │   └── unstable -> jammy
│   ├── EFI
│   │   └── boot
│   ├── install
│   ├── md5sum.txt
│   ├── nocloud
│   │   ├── dcloud.tar.gz.0
│   │   ├── dcloud.tar.gz.1
│   │   ├── meta-data
│   │   └── user-data
│   ├── pool
│   │   ├── main
│   │   └── restricted
│   └── ubuntu -> .
├── ubuntu-autoinstall-generator.sh
└── user-data.example
 
16 directories, 35 files


4. 打包
打包脚本：

xorriso -as mkisofs -r \
-V 'Ubuntu-Server 22.04.1 LTS amd64' \
-o ./ubuntu-22.04.1-dcloud-v000.iso \
--modification-date='2022080916483300' \
--grub2-mbr ./BOOT/1-Boot-NoEmul.img \
--protective-msdos-label \
-partition_cyl_align off \
-partition_offset 16 \
--mbr-force-bootable \
-append_partition 2 28732ac11ff8d211ba4b00a0c93ec93b ./BOOT/2-Boot-NoEmul.img \
-appended_part_as_gpt \
-iso_mbr_part_type a2a0d0ebe5b9334487c068b6b72699c7 \
-c '/boot.catalog' \
-b '/boot/grub/i386-pc/eltorito.img' \
-no-emul-boot \
-boot-load-size 4 \
-boot-info-table \
--grub2-boot-info \
-eltorito-alt-boot \
-e '--interval:appended_partition_2_start_717863s_size_8496d:all::' \
-no-emul-boot \
-boot-load-size 8496 \
./ubuntu-autoinstall-generator/tmpdir/


执行上面脚本进行打包

备注：命令可参考源于：xorriso -indev xxxxxx.iso -report_el_torito as_mkisofs











安装过程
1) 选择 Try install ....

2) 配置网络 （网关和DNS 都必须填， 不能留空）



3） 配置磁盘 （确保 “/” 的挂载点至少有100G以上的空间容量）



4） 等待安装结束后reboot

