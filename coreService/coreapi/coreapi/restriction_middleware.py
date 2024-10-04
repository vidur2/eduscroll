from django.http import HttpResponseForbidden

class HostRestrictionMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        allowed_hosts = {'localhost', 'edutok-textbook-1', 'edutok-jit-1', 'edutok-coreservice-1'}
        current_host = request.get_host()
        current_path = request.path

        # Specify the path you want to restrict
        restricted_paths = {'/embed'}
        print(f"current_host: {current_host}")
        if current_path in restricted_paths and current_host.split(":")[0] not in allowed_hosts:
            return HttpResponseForbidden("Forbidden: Host not allowed for this API route.")

        response = self.get_response(request)
        return response
