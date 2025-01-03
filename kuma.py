#https://github.com/lucasheld/uptime-kuma-api/
#https://uptime-kuma-api.readthedocs.io/en/latest/index.html

from uptime_kuma_api import UptimeKumaApi, MonitorType

# 文件路径
domains_file = "domains.txt"

# Uptime Kuma 服务器地址和登录信息
uptime_kuma_url = "http://10.10.20.23"
username = "urlitor"
password = "wnPu1scWVKjm"

# 从文件读取域名
def read_domains_from_file(file_path):
    try:
        with open(file_path, "r") as file:
            lines = file.readlines()
            # 去除空行和换行符
            domains = [line.strip() for line in lines if line.strip()]
            return domains
    except FileNotFoundError:
        print(f"Error: File '{file_path}' not found.")
        return []

# 批量添加监控
def batch_add_monitors():
    try:
        # 初始化 API 实例
        api = UptimeKumaApi(uptime_kuma_url)

        # 登录
        api.login(username, password)
        print("Logged in successfully!")

        # 读取 URL 列表
        domains = read_domains_from_file(domains_file)
        if not domains:
            print("No domains found in the file.")
            return

        # 循环添加监控
        for index, domain in enumerate(domains, start=1):
            try:
                result = api.add_monitor(
                        type=MonitorType.KEYWORD, 
                        name=domain,
                        parent=1,
                        interval=300,
                        maxretries=1,
                        ignoreTls=1,
                        keyword="OK",
                        maxredirects=5,
                        accepted_statuscodes=['200-299', '300-399'],
                        expiryNotification=1,
                        timeout='10',
                        url="https://" + domain +"/monitor.txt"
                )
                print(f"Monitor added for {domain}: {result}")
            except Exception as e:
                print(f"Failed to add monitor for {domain}: {str(e)}")

        # 登出
        api.logout()
        print("Logged out successfully!")

    except Exception as e:
        print(f"Error: {str(e)}")

# 执行批量添加
if __name__ == "__main__":
    batch_add_monitors()
