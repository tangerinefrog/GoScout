const statuses = [
    'created',
    'graded',
    'ignored',
    'applied',
]

async function getJobs() {
    const resp = await fetch(`/api/jobs`);
    if (!resp.ok) {
        showErrorToast('Jobs fetch failed');
        return [];
    }

    const jobs = await resp.json();
    if (!jobs) {
        return [];
    }

    return jobs;    
}

async function updateJob(id, field, value) {
    const reqBody = {
        [field]: value 
    };

    const resp = await fetch(`/api/jobs/${id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(reqBody),
    });

    return resp.ok;
}