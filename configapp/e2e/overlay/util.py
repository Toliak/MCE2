import subprocess
import sys


def check_output_with_live_echo(
    args,
    executable=None,
    stdin=None,
    preexec_fn=None,
    close_fds=True,
    shell=False,
    cwd=None,
    env=None,
    universal_newlines=None,
    startupinfo=None,
    creationflags=0,
    restore_signals=True,
    start_new_session=False,
    pass_fds=(),
    *,
    user=None,
    group=None,
    extra_groups=None,
    encoding=None,
    errors=None,
    umask=-1,
    pipesize=-1,
):
    """
    Runs a command like subprocess.check_output, but:
    1. Merges stderr into stdout
    2. Captures combined output into a buffer (returned)
    3. Echoes output live to console stdout
    """
    process = subprocess.Popen(
        args,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        bufsize=1,
        text=True,
        ##
        executable=executable,
        stdin=stdin,
        preexec_fn=preexec_fn,
        close_fds=close_fds,
        shell=shell,
        cwd=cwd,
        env=env,
        universal_newlines=universal_newlines,
        startupinfo=startupinfo,
        creationflags=creationflags,
        restore_signals=restore_signals,
        start_new_session=start_new_session,
        pass_fds=pass_fds,
        user=user,
        group=group,
        extra_groups=extra_groups,
        encoding=encoding,
        errors=errors,
        umask=umask,
        pipesize=pipesize,
    )

    output_chunks = []

    try:
        for line in process.stdout:
            sys.stdout.write(line)
            
            output_chunks.append(line)
    finally:
        process.stdout.close()

    sys.stdout.flush()
    return_code = process.wait()
    output = "".join(output_chunks)

    if return_code != 0:
        raise subprocess.CalledProcessError(return_code, args, output=output)

    return output
