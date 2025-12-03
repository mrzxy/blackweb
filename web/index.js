// 全局变量
let trade_data = [];
let visibleData = [];
let isLoading = false;
let currentTimezone = 'EST'; // 默认时区
let lastCreationEpochSec = 0; // 增量拉取基准（秒，UTC当日）
let pollTimerId = null; // 定时器ID
let currentOffset = 0; // 当前加载的偏移量
let hasMoreData = true; // 是否还有更多数据

// DOM元素
let tableBody, filterBtn, filterModal, saveBtn, applyBtn, resetBtn, cancelBtn, searchBox, tableContainer, loading, loader;

// 初始化函数
function init() {
    // 获取DOM元素
    tableBody = document.getElementById('table-body');
    filterBtn = document.getElementById('filter-btn');
    filterModal = document.getElementById('filter-modal');
    saveBtn = filterModal.querySelector('.save');
    applyBtn = filterModal.querySelector('.apply');
    resetBtn = filterModal.querySelector('.reset');
    cancelBtn = filterModal.querySelector('.cancel');
    searchBox = document.getElementById('search-box');
    tableContainer = document.getElementById('table-container');
    loading = document.getElementById('loading');
    loader = document.getElementById('loader');
    
    // 加载时区设置
    loadTimezoneSettings();
    
    // 绑定事件
    bindEvents();
    
    // 加载已保存的筛选配置，如果没有保存的配置则使用默认状态
    loadSavedFilters();
    
    applyFilterSettings()

    startIncrementalPolling();
}

// 绑定事件
function bindEvents() {
    // 无限滚动
    tableContainer.addEventListener('scroll', () => {
        const { scrollTop, scrollHeight, clientHeight } = tableContainer;
        
        if (scrollTop + clientHeight >= scrollHeight - 100) {
            loadMoreData();
        }
    });
    
    // 搜索功能
    searchBox.addEventListener('input', function() {
        const searchTerm = this.value.toLowerCase();
        
        if (searchTerm === '') {
            visibleData = [...trade_data];
            renderTable(visibleData);
            return;
        }
        
        const filteredData = visibleData.filter(item => 
            (item.Symbol && item.Symbol.toLowerCase().includes(searchTerm)) || 
            (item.Details && item.Details.toLowerCase().includes(searchTerm)) ||
            (item.CP && item.CP.toLowerCase().includes(searchTerm))
        );
        
        renderTable(filteredData);
    });
    
    // 筛选设置
    filterBtn.addEventListener('click', () => {
        filterModal.style.display = 'flex';
    });
    
    // 保存筛选设置
    saveBtn.addEventListener('click', () => {
        applyFilterSettings();
        saveFilterSettings();
        filterModal.style.display = 'none';
    });
    
    // 应用筛选设置
    applyBtn.addEventListener('click', () => {
        applyFilterSettings();
        filterModal.style.display = 'none';
    });
    
    // 重置筛选设置
    resetBtn.addEventListener('click', () => {
        // 清除本地存储的筛选设置
        localStorage.removeItem('optionFilters');
        
        // 重置为默认状态
        setDefaultFilterState();
        
        // 重新渲染表格
        visibleData = [...trade_data];
        renderTable(visibleData);
        
        console.log('筛选设置已重置为默认状态');
    });
    
    // 取消筛选设置
    cancelBtn.addEventListener('click', () => {
        filterModal.style.display = 'none';
    });
    
    // 点击模态框外部关闭
    window.addEventListener('click', (e) => {
        if (e.target === filterModal) {
            filterModal.style.display = 'none';
        }
    });
}

// 加载交易数据
function loadTradeData(queryConfig, is_incremental, is_load_more) {

	let req = {
		"limit": 100, // 默认限制
		"offset": 0
	};

    req = { ...req, ...queryConfig };
	
	post('/api/option-trades', req, function(resp) {
		if (resp.code === 200) {
			if (is_incremental) {
				// 增量更新（新数据插入到顶部）
				if (Array.isArray(resp.data) && resp.data.length > 0) {
					for (let i = resp.data.length - 1; i >= 0; i--) {
						trade_data.unshift(resp.data[i]);
					}
					visibleData = [...trade_data];
					renderTable(visibleData);
				}
			} else if (is_load_more) {
				// 加载更多（追加到底部）
				if (Array.isArray(resp.data) && resp.data.length > 0) {
					trade_data = trade_data.concat(resp.data);
					visibleData = [...trade_data];
					renderTable(visibleData);
					currentOffset += resp.data.length; // 更新偏移量
					hasMoreData = resp.data.length === req.limit; // 如果返回数据少于limit，说明没有更多数据了
				} else {
					hasMoreData = false;
				}
			} else {
				// 初始加载
				trade_data = resp.data;
				visibleData = [...trade_data];
				currentOffset = resp.data.length; // 更新偏移量
				hasMoreData = resp.data.length === req.limit;
				renderTable(visibleData);
			}
		} else {
			console.error('加载数据失败:', resp.message);
		}
	});
}

// 将筛选配置转换为 QueryRequest 格式
function convertFiltersToQueryRequest(filters) {
	const queryConfig = {
        BidAsk: [],
        PreValue: [],
    };
	
	// 遍历筛选组
	for (const groupName in filters) {
		const groupOptions = filters[groupName];
		
		groupOptions.forEach(option => {
			const optionValue =  option.value 

            if (option.name === 'pre_value') {
                queryConfig.PreValue.push(optionValue);
                return 
            }
			
			// 根据选项名称映射到 QueryRequest 字段
			switch (optionValue) {
				case 'put':
				case 'call':
					if (!queryConfig.optionType) queryConfig.optionType = [];
					queryConfig.optionType.push(optionValue.toUpperCase());
					break;
					
				case 'yellow':
				case 'white':
				case 'magenta':
					if (!queryConfig.flowColor) queryConfig.flowColor = [];
					queryConfig.flowColor.push(optionValue.toUpperCase());
					break;
					
				case 'etf':
				case 'stock':
					if (!queryConfig.securityType) queryConfig.securityType = [];
					queryConfig.securityType.push(optionValue.toUpperCase());
					break;
					
				case 'AA':
				case 'BB':
				case 'A':
				case 'B':
					queryConfig.BidAsk.push(optionValue); 
					break;
                case 'lt$0.75':
                    queryConfig.marketCapAbove750B = true;
					break;
				case 'in-the-money':	
                    queryConfig.inTheMoney = true;
                    break;
				case 'in-the-money':	
                    queryConfig.inTheMoney = true;
                    break;
                case 'out-the-money':
                    queryConfig.outTheMoney = true;
                    break;
                case 'sweep-only':
                    queryConfig.weepOnly = true;
                    break;
                case 'weekly-only':
                    queryConfig.weeklyOnly = true;
                    break;
                case 'earnings':
                    queryConfig.earnings = true;
                    break;
                case 'unusual':
                    queryConfig.unusual = true;
                    break;
       


				case 'consumer-discretionary':
				case 'industrials':
				case 'information-technology':
				case 'real-estate':
				case 'health-care':
				case 'energy':
				case 'financials':
				case 'materials':
				case 'consumer-staples':
				case 'communication-services':
				case 'utilities':
					// 行业分类，可以添加到 sector 字段（如果后端支持）
					if (!queryConfig.sector) queryConfig.sector = [];
					queryConfig.sector.push(optionValue.replace(/-/g, ' ').replace(/\b\w/g, l => l.toUpperCase()));
					break;
			}
		});
	}
	
	return queryConfig;
}

let last_creation_date = ""
// 渲染表格
function renderTable(data) {
    if (!tableBody) return;
    
    tableBody.innerHTML = '';

    let tmp_last_creation_date = ""
    
    data.forEach(item => {
        const row = document.createElement('tr');

        if (last_creation_date == "" || parseInt(tmp_last_creation_date)  < parseInt(item.CreationDate)) {
            last_creation_date = item.CreationDate
        }
        
        // 根据 color 字段设置行的配色方案
        const color = item.Color || '';
        if (color === 'MAGENTA') {
            row.style.color = 'magenta';
            row.style.backgroundColor = 'transparent';
        } else if (color === 'YELLOW') {
            row.style.color = 'yellow';
            row.style.backgroundColor = 'transparent';
        } else {
            // 默认配色
            row.style.color = 'white';
            row.style.backgroundColor = 'rgb(33, 42, 51)';
        }
        
        row.innerHTML = `
            <td class="time">${item.Time}</td>
            <td class="symbol">${item.Symbol || ''}</td>
            <td>${item.Exp}</td>
            <td class="strike">${formatStrike(item.Strike)}</td>
            <td class="cp cp-${(item.CP || '').toLowerCase()}" style="color: ${(item.CP || '').toUpperCase() === 'CALL' ? 'lightgreen' : 'red'}">${item.CP || ''}</td>
            <td class="spot">${formatSpot(item.Spot)}</td>
            <td class="details">${item.Details || ''}</td>
            <td><span class="trade-type ${(item.Type || '').toLowerCase()}">${item.Type || ''}</span></td>
            <td class="value"><span class="value-highlight">${formatValue(item.Value)}</span></td>
            <td class="iv"><span class="iv-highlight">${item.Iv || '0'}</span></td>
        `;
        tableBody.appendChild(row);
    });
    if (last_creation_date != "" && parseInt(tmp_last_creation_date) > parseInt(last_creation_date)) {
        last_creation_date = tmp_last_creation_date
    }
    
    const headerTitle = document.querySelector('.header-title');
    if (headerTitle) {
        headerTitle.textContent = `OPTIONS (${data.length}) > ALL`;
    }
}

// 格式化时间
function formatTime(timestamp) {
    if (!timestamp) return '';
    const date = new Date(timestamp);
    return date.toLocaleTimeString('en-US', { 
        hour: '2-digit', 
        minute: '2-digit', 
        second: '2-digit',
        hour12: false 
    });
}

// 格式化Unix时间戳（毫秒）为HH:MM:SS格式
function formatUnixTime(unixTimestamp) {
    if (!unixTimestamp) return '';
    
    // 将毫秒时间戳转换为Date对象
    const date = new Date(parseInt(unixTimestamp));
    
    // 检查日期是否有效
    if (isNaN(date.getTime())) return '';
    
    // 格式化为HH:MM:SS
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    const seconds = date.getSeconds().toString().padStart(2, '0');
    
    return `${hours}:${minutes}:${seconds}`;
}



// 将多种输入格式的 Exp 统一为 MM/DD/YY
function formatExp(exp) {
    if (!exp) return '';
    // 如果可被 Date 解析（例如 ISO 字符串）
    const tryDate = new Date(parseInt(exp));
    if (!isNaN(tryDate.getTime())) {
        const mm = String(tryDate.getMonth() + 1).padStart(2, '0');
        const dd = String(tryDate.getDate()).padStart(2, '0');
        const yy = String(tryDate.getFullYear()).slice(-2);
        return `${mm}/${dd}/${yy}`;
    }
    // 手动解析常见的分隔格式 MM/DD/YYYY, YYYY-MM-DD, M/D/YY 等
    const m = String(exp).trim();
    // 1) 形如 YYYY-MM-DD 或 YYYY/M/D
    let match = m.match(/^\s*(\d{4})[-\/.](\d{1,2})[-\/.](\d{1,2})\s*$/);
    if (match) {
        const mm = match[2].padStart(2, '0');
        const dd = match[3].padStart(2, '0');
        const yy = match[1].slice(-2);
        return `${mm}/${dd}/${yy}`;
    }
    // 2) 形如 MM/DD/YYYY 或 M/D/YYYY
    match = m.match(/^\s*(\d{1,2})[\/-](\d{1,2})[\/-](\d{4})\s*$/);
    if (match) {
        const mm = match[1].padStart(2, '0');
        const dd = match[2].padStart(2, '0');
        const yy = match[3].slice(-2);
        return `${mm}/${dd}/${yy}`;
    }
    // 3) 形如 MM/DD/YY 或 M/D/YY，补零
    match = m.match(/^\s*(\d{1,2})[\/-](\d{1,2})[\/-](\d{2})\s*$/);
    if (match) {
        const mm = match[1].padStart(2, '0');
        const dd = match[2].padStart(2, '0');
        const yy = match[3];
        return `${mm}/${dd}/${yy}`;
    }
    // 其他未知格式，原样返回
    return m;
}

// 格式化行权价
function formatStrike(strike) {
    if (!strike) return '$0';
    const num = parseFloat(strike);
    return `$${num.toFixed(2)}`;
}

// 格式化现价
function formatSpot(spot) {
    if (!spot) return '$0';
    const num = parseFloat(spot);
    return `$${num.toFixed(2)}`;
}

// 格式化权利金
function formatPremium(premium) {
    if (!premium) return '$0';
    const num = parseFloat(premium);
    if (num >= 1000000) {
        return `$${(num / 1000000).toFixed(2)}M`;
    } else if (num >= 1000) {
        return `$${(num / 1000).toFixed(1)}K`;
    } else {
        return `$${num.toFixed(0)}`;
    }
}

// 格式化价值
function formatValue(value) {
    if (!value) return '$0';
    const num = parseFloat(value);
    if (num >= 1000000) {
        return `$${(num / 1000000).toFixed(2)}M`;
    } else if (num >= 1000) {
        return `$${(num / 1000).toFixed(1)}K`;
    } else {
        return `$${num.toFixed(0)}`;
    }
}

// 加载更多数据
function loadMoreData() {
    if (isLoading || !hasMoreData) return;
    
    isLoading = true;
    if (loading) loading.style.display = 'block';
    if (loader) loader.style.display = 'block';
    
    // 获取当前筛选配置
    let savedFilters = getSelectedFilters();
    let queryConfig = {
        offset: currentOffset,
        limit: 100
    };
    
    if (savedFilters) {
        queryConfig = { ...queryConfig, ...convertFiltersToQueryRequest(savedFilters) };
    }
    
    // 调用 loadTradeData 并传递 is_load_more=true
    loadTradeData(queryConfig, false, true);
    
    // 注意：isLoading 的重置会在 loadTradeData 的回调中处理
    // 这里需要在 loadTradeData 完成后重置状态
    setTimeout(() => {
        isLoading = false;
        if (loading) loading.style.display = 'none';
        if (loader) loader.style.display = 'none';
    }, 500);
}

// 应用筛选设置
function applyFilterSettings() {
    let queryConfig = {}
    let savedFilters = getSelectedFilters()
    if (savedFilters) {
        queryConfig = convertFiltersToQueryRequest(savedFilters);
	}

    // 重置偏移量和状态
    currentOffset = 0;
    hasMoreData = true;
    
    loadTradeData(queryConfig, false, false);
}

function getSelectedFilters() {
    const filters = {};
    
    // 获取所有筛选状态
    document.querySelectorAll('.filter-group').forEach(group => {
        const groupName = group.querySelector('h3').textContent.trim();
        const options = [];
        
        group.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
            if (checkbox.checked) {
                options.push({
                    name: checkbox.name,
                    value: checkbox.value
                });
            }
        });
        
        filters[groupName] = options;
    });
    return filters;
}

// 保存筛选设置到本地存储
function saveFilterSettings() {
    const filters = getSelectedFilters();
    // 保存到本地存储
    localStorage.setItem('optionFilters', JSON.stringify(filters));
    alert('筛选设置已保存！');
}

// 加载保存的筛选设置
function loadSavedFilters() {
	const savedFilters = localStorage.getItem('optionFilters');
	
	if (savedFilters) {
		// 如果有保存的配置，应用已保存的设置
		const filters = JSON.parse(savedFilters);
		
		// 获取所有复选框
		const allCheckboxes = document.querySelectorAll('#filter-modal input[type="checkbox"]');
		
		// 先取消所有选项的选中状态
		allCheckboxes.forEach(checkbox => {
			checkbox.checked = false;
		});
		
		// 应用已保存的设置到UI
		for (const groupName in filters) {
			const groupOptions = filters[groupName];
			groupOptions.forEach(option => {
				const checkbox = document.querySelector(`input[name="${option.name}"][value="${option.value}"]`);
				if (checkbox) {
					checkbox.checked = true;
				}
			});
		}
		
		console.log('已加载保存的筛选配置');
	} else {
		// 如果没有保存的配置，使用默认状态
		setDefaultFilterState();
		console.log('没有保存的筛选配置，使用默认状态');
	}

}

// 设置默认筛选状态
function setDefaultFilterState() {
	// 定义默认应该选中的选项
	const defaultSelectedOptions = [
		'put',
		'call', 
		'yellow',
		'white',
        'magenta',
		'AA',
		'A',
		'etf',
		'stock',
		'consumer-discretionary',
		'industrials',
		'information-technology',
		'real-estate',
		'health-care',
		'energy',
		'financials',
		'materials',
		'consumer-staples',
		'communication-services',
		'utilities'
	];
	
	// 获取所有复选框
	const allCheckboxes = document.querySelectorAll('#filter-modal input[type="checkbox"]');
	
	// 先取消所有选项的选中状态
	allCheckboxes.forEach(checkbox => {
		checkbox.checked = false;
	});
	
	// 只选中默认指定的选项
	allCheckboxes.forEach(checkbox => {
		const checkboxValue = checkbox.value ;
		if (defaultSelectedOptions.includes(checkboxValue)) {
			checkbox.checked = true;
		}
	});
	

	
	console.log('筛选设置已设置为默认状态：指定选项打钩，Above Ask 固定选中');
}

// 更新时间
function updateTime() {
    const now = new Date();
    const options = { 
        hour: '2-digit', 
        minute: '2-digit', 
        second: '2-digit',
        hour12: false 
    };
    const timeString = now.toLocaleTimeString('en-US', options);
    
    const statusBar = document.querySelector('.status-bar > div:first-child');
    if (statusBar) {
        statusBar.textContent = `更新于: ${now.toLocaleDateString()} ${timeString} ${currentTimezone}`;
    }
}

// 切换时区
function changeTimezone(newTimezone) {
    currentTimezone = newTimezone;
    
    // 重新渲染表格以更新时间显示
    if (visibleData.length > 0) {
        renderTable(visibleData);
    }
    
    // 更新时间显示
    updateTime();
    
    // 保存时区设置到本地存储
    localStorage.setItem('preferredTimezone', currentTimezone);
    
    console.log(`时区已切换到: ${currentTimezone}`);
}

// 获取支持的时区列表
function getSupportedTimezones() {
    return [
        { code: 'EST', name: 'Eastern Standard Time (UTC-5)' },
        { code: 'PST', name: 'Pacific Standard Time (UTC-8)' },
        { code: 'CST', name: 'Central Standard Time (UTC-6)' },
        { code: 'MST', name: 'Mountain Standard Time (UTC-7)' },
        { code: 'UTC', name: 'UTC Time' }
    ];
}

// 加载时区设置
function loadTimezoneSettings() {
    const savedTimezone = localStorage.getItem('preferredTimezone');
    if (savedTimezone) {
        currentTimezone = savedTimezone;
        console.log(`已加载保存的时区设置: ${currentTimezone}`);
    }
}

// AJAX请求函数
function post(url, data, callback) {
    $.ajax({  
        type: "post",  
        url: url,  
        async: false, // 使用同步方式  
        data: JSON.stringify(data),  
        contentType: "application/json; charset=utf-8",  
        dataType: "json",  
        success: callback
    });  
}



// 增量拉取
function fetchIncrementalData() {
    console.log("fetchIncrementalData", last_creation_date)
    if (last_creation_date == "") {
        return
    }
    let savedFilters = getSelectedFilters()
    let queryConfig = {
    }
    if (savedFilters) {
        queryConfig = convertFiltersToQueryRequest(savedFilters);
    }
    queryConfig.lastCreationate = last_creation_date
    

    loadTradeData(queryConfig, true, false);
}

// 启动/重启轮询
function startIncrementalPolling() {
	if (pollTimerId) {
		clearInterval(pollTimerId);
		pollTimerId = null;
	}

	pollTimerId = setInterval(fetchIncrementalData, 2000);
}

// 页面加载完成后初始化
$(function(){
    // 等待DOM加载完成
    $(document).ready(function() {
        init();
        // 启动时间更新
        updateTime();
        setInterval(updateTime, 1000);
    });
});