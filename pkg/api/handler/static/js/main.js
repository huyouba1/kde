// 页面加载完成后执行
document.addEventListener('DOMContentLoaded', function() {
    // 初始化页面数据
    initializeDashboard();

    // 设置导航栏活动状态
    setActiveNavItem();
});

// 初始化仪表板数据
async function initializeDashboard() {
    try {
        // 从API获取最新数据
        const response = await fetch('/api/v1/dashboard');
        if (!response.ok) {
            throw new Error('获取数据失败');
        }

        const data = await response.json();
        updateDashboard(data);
    } catch (error) {
        console.error('加载仪表板数据失败:', error);
        showError('无法加载仪表板数据');
    }
}

// 更新仪表板显示
function updateDashboard(data) {
    // 更新统计卡片
    updateStatCard('cluster-count', data.clusterCount);
    updateStatCard('deployment-count', data.deploymentCount);
    updateStatCard('system-status', data.systemStatus);
}

// 更新统计卡片数值
function updateStatCard(id, value) {
    const element = document.querySelector(`#${id} .stat-value`);
    if (element) {
        element.textContent = value;
    }
}

// 设置当前活动的导航项
function setActiveNavItem() {
    const currentPath = window.location.pathname;
    const navLinks = document.querySelectorAll('.nav a');

    navLinks.forEach(link => {
        if (link.getAttribute('href') === currentPath) {
            link.classList.add('active');
        }
    });
}

// 显示错误消息
function showError(message) {
    // TODO: 实现错误提示UI组件
    console.error(message);
}